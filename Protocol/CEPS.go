package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	netpack "MPC/Netpackage"
	party "MPC/Party"
	parsing "MPC/parsing"
	"fmt"
	"math"
	"strconv"
	"sync"
)

type Ceps struct {
	config        *config.Config
	peer          party.IPeer
	shamir        *ShamirSecretSharing
	cMessages     chan netpack.Share
	rShares       rShares
	degree        int
	fShares       fShares
	subscribeMap  subscribeMap
	listOfRandoms []randoms
	matrix        [][]int64
}

type subscribeMap struct {
	m    map[netpack.ShareIdentifier]chan int
	lock sync.Mutex
}

func (subM *subscribeMap) ping(iden netpack.ShareIdentifier) {
	subM.lock.Lock()
	c := subM.m[iden]
	subM.lock.Unlock()
	if c != nil {
		c <- 1
	}

}

type fShares struct {
	finalShares []netpack.Share
	mu          sync.Mutex
}

type rShares struct {
	receivedShares map[netpack.ShareIdentifier]*netpack.Share
	mu             sync.Mutex
}

func mkProtocol(config *config.Config, field field.Field, peer party.IPeer) *Ceps {
	prot := new(Ceps)
	prot.cMessages = make(chan netpack.Share)
	prot.config = config
	prot.peer = peer
	prot.degree = int(math.Ceil(prot.config.ConstantConfig.NumberOfParties/2) - 1)
	prot.shamir = makeShamirSecretSharing(config.VariableConfig.Secret, field, prot.degree)
	prot.rShares = rShares{receivedShares: make(map[netpack.ShareIdentifier]*netpack.Share)}
	prot.subscribeMap = subscribeMap{m: make(map[netpack.ShareIdentifier]chan int)}
	prot.matrix = prot.createMatrix()
	return prot
}

func (prot *Ceps) run() int64 {
	var wg sync.WaitGroup
	//Start peer
	wg.Add(1)
	prot.peer.StartPeer(prot.cMessages, &wg)

	//wait for network
	wg.Wait()

	//Start receiving messages from the network

	go prot.receive()
	//Convert string expression into instruction list
	exp := prot.config.ConstantConfig.Expression
	astExp := parsing.ParseExpression(exp)
	instructionTree, err := parsing.ConvertAstToTree(astExp)
	if err != nil {
		//TODO maybe shut down peer?
		println(err.Error())
		return 0
	}
	prot.runPreprocess(instructionTree.CountMults())
	//Send input shares
	toSendIdentifier := netpack.ShareIdentifier{Ins: "p" + strconv.Itoa(int(prot.config.VariableConfig.PartyNr)), PartyNr: int(prot.config.VariableConfig.PartyNr)}
	shares := prot.shamir.makeShares(int64(prot.config.ConstantConfig.NumberOfParties), toSendIdentifier)
	prot.handleShare(shares)

	//Do instructions
	prot.calculateInstruction(instructionTree)

	//output reconstruction
	res := prot.outputReconstruction(instructionTree.GetResultName())
	return res
}

func (prot *Ceps) calculateInstruction(instructionTree *parsing.InstructionTree) {

	if instructionTree.Left != nil {
		go prot.calculateInstruction(instructionTree.Left)
	}
	if instructionTree.Right != nil {
		go prot.calculateInstruction(instructionTree.Right)
	}
	switch ins := instructionTree.Instruction.(type) {
	case *parsing.AddInstruction:
		prot.add(ins)
	case *parsing.MultInstruction:
		prot.multiply(ins)
	case *parsing.ScalarInstruction:
		prot.scalar(ins)
	default:
		fmt.Printf("Unknown instruction %v", ins)
	}
	// return
}

func (prot *Ceps) receive() {
	for {
		message := <-prot.cMessages
		fmt.Printf("Party %v got share {%v,%v}\n", prot.config.VariableConfig.PartyNr, message.Identifier, message.Value)
		if string(message.Identifier.Ins) == "o" {
			prot.fShares.mu.Lock()
			prot.fShares.finalShares = append(prot.fShares.finalShares, message)
			prot.fShares.mu.Unlock()
		} else {
			prot.rShares.mu.Lock()
			prot.rShares.receivedShares[message.Identifier] = &message
			prot.rShares.mu.Unlock()
		}
		go prot.subscribeMap.ping(message.Identifier)

	}
}

func (prot *Ceps) waitForShares(needToWaitOn []netpack.ShareIdentifier) {
	var waitChannels []chan int
	prot.subscribeMap.lock.Lock()
	prot.rShares.mu.Lock()
	for _, s := range needToWaitOn {
		if prot.rShares.receivedShares[s] == nil {
			c := make(chan int)
			prot.subscribeMap.m[s] = c
			waitChannels = append(waitChannels, c)
		}
	}
	prot.rShares.mu.Unlock()
	prot.subscribeMap.lock.Unlock()
	for _, c := range waitChannels {
		<-c
	}
}

func (prot *Ceps) addResultShare(insResult string, value int64) {
	resultShare := &netpack.Share{}
	resultShare.Value = value
	resultShare.Identifier = netpack.ShareIdentifier{Ins: insResult, PartyNr: int(prot.config.VariableConfig.PartyNr)}
	prot.rShares.mu.Lock()
	prot.rShares.receivedShares[resultShare.Identifier] = resultShare
	prot.rShares.mu.Unlock()
	go prot.subscribeMap.ping(resultShare.Identifier)
}

func (prot *Ceps) handleShare(shares []netpack.Share) {
	fmt.Printf("Party %v send %v\n", prot.config.VariableConfig.PartyNr, shares)
	go prot.peer.SendShares(shares)
}

func (prot *Ceps) createWaitShareIdentifier(ins string) netpack.ShareIdentifier {
	var partyNr int
	if string(ins[0]) == "p" {
		partyNr, _ = strconv.Atoi(string(ins[1:]))

	} else {
		partyNr = int(prot.config.VariableConfig.PartyNr)
	}
	return netpack.ShareIdentifier{Ins: ins, PartyNr: partyNr}

}

func (prot *Ceps) add(ins *parsing.AddInstruction) {
	leftIden := prot.createWaitShareIdentifier(ins.Left)
	rightIden := prot.createWaitShareIdentifier(ins.Right)
	prot.waitForShares([]netpack.ShareIdentifier{leftIden, rightIden})
	prot.rShares.mu.Lock()
	leftVal := prot.rShares.receivedShares[leftIden].Value
	rightVal := prot.rShares.receivedShares[rightIden].Value
	prot.rShares.mu.Unlock()
	value := prot.shamir.field.Add(leftVal, rightVal)
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) multiply(ins *parsing.MultInstruction) {
	leftIden := prot.createWaitShareIdentifier(ins.Left)
	rightIden := prot.createWaitShareIdentifier(ins.Right)
	prot.waitForShares([]netpack.ShareIdentifier{leftIden, rightIden})
	prot.rShares.mu.Lock()
	a := prot.rShares.receivedShares[leftIden].Value
	b := prot.rShares.receivedShares[rightIden].Value
	prot.rShares.mu.Unlock()
	ab2t := prot.shamir.field.Multiply(a, b)
	if ins.Num > len(prot.listOfRandoms) {
		fmt.Printf("Impossible - party %v did not have enough r-values for mult %v\n", prot.config.VariableConfig.PartyNr, ins.Num)
		fmt.Printf("party %v r-values: %v\n", prot.config.VariableConfig.PartyNr, prot.listOfRandoms)
		return
	}
	rPair := prot.listOfRandoms[ins.Num-1]
	abMinusrShare := prot.shamir.field.Minus(ab2t, rPair.r2t)
	//Send to party ins.Num mod n to distribute load
	toSendIdentifier := netpack.ShareIdentifier{Ins: "m" + strconv.Itoa(ins.Num), PartyNr: int(prot.config.VariableConfig.PartyNr)}
	sendTo := ((ins.Num - 1) % int(prot.config.ConstantConfig.NumberOfParties)) + 1
	prot.peer.SendShare(netpack.Share{Value: abMinusrShare, Identifier: toSendIdentifier}, sendTo)

	abrIden := netpack.ShareIdentifier{Ins: "ab-r" + strconv.Itoa(ins.Num), PartyNr: sendTo}
	if int(prot.config.VariableConfig.PartyNr) == sendTo {
		//If I am receiver then I need to receive and compute ab-r
		var multiplicationIdentifiers []netpack.ShareIdentifier
		for i := 1; i <= int(prot.config.ConstantConfig.NumberOfParties); i++ {
			multiplicationIdentifiers = append(multiplicationIdentifiers, netpack.ShareIdentifier{Ins: "m" + strconv.Itoa(ins.Num), PartyNr: i})
		}
		prot.waitForShares(multiplicationIdentifiers)
		var sharesForLagrange []netpack.Share
		for _, i := range multiplicationIdentifiers {
			sharesForLagrange = append(sharesForLagrange, *prot.rShares.receivedShares[i])
		}
		value, _ := prot.shamir.lagrangeInterpolation(sharesForLagrange, int(prot.config.ConstantConfig.NumberOfParties-1))
		//Then share ab-r as the constant polynomial
		abrShare := netpack.Share{Value: value, Identifier: abrIden}
		var abrShares []netpack.Share
		for i := 0; i < int(prot.config.ConstantConfig.NumberOfParties); i++ {
			abrShares = append(abrShares, abrShare)
		}
		prot.handleShare(abrShares)
	}
	//Wait for ab-r
	prot.waitForShares([]netpack.ShareIdentifier{abrIden})
	//Each party computes the share of ab_t
	prot.rShares.mu.Lock()
	abMinusr := prot.rShares.receivedShares[abrIden].Value
	prot.rShares.mu.Unlock()
	insResultValue := prot.shamir.field.Add(rPair.r1t, abMinusr)
	prot.addResultShare(ins.Result, insResultValue)
}

func (prot *Ceps) scalar(ins *parsing.ScalarInstruction) {
	rightIden := prot.createWaitShareIdentifier(ins.Variable)
	prot.waitForShares([]netpack.ShareIdentifier{rightIden})
	prot.rShares.mu.Lock()
	varValue := prot.rShares.receivedShares[rightIden].Value
	prot.rShares.mu.Unlock()
	scalar := prot.shamir.field.Convert(int64(ins.Scalar))
	value := prot.shamir.field.Multiply(scalar, varValue)
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) outputReconstruction(finalResult string) int64 {
	resIden := prot.createWaitShareIdentifier(finalResult)
	//Send out result share
	prot.rShares.mu.Lock()
	resultShare := prot.rShares.receivedShares[resIden]
	prot.rShares.mu.Unlock()
	resultShare.Identifier = netpack.ShareIdentifier{Ins: "o", PartyNr: int(prot.config.VariableConfig.PartyNr)}
	prot.peer.SendFinal(*resultShare)
	prot.fShares.mu.Lock()
	prot.fShares.finalShares = append(prot.fShares.finalShares, *resultShare)
	prot.fShares.mu.Unlock()
	shares := prot.waitForFinalShares()
	fmt.Printf("Party %v is calling lagrange with degree: %v, and shares: %v\n", prot.config.VariableConfig.PartyNr, prot.degree, shares)
	result, err := prot.shamir.lagrangeInterpolation(shares, prot.degree)
	if err != nil {
		println(err)
		return 0
	}
	return result
}

func (prot *Ceps) waitForFinalShares() []netpack.Share {
	//Could be made with observer pattern, register an observer when this is called, an observer
	//could just be a chan int, each time you get a new package you do notify which sends a signal on all (multiple) channels for each wait method
	for {
		prot.fShares.mu.Lock()
		if len(prot.fShares.finalShares) > prot.degree {
			break
		}
		prot.fShares.mu.Unlock()
	}
	return prot.fShares.finalShares

}
