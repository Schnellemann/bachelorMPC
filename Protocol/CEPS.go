package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	netpack "MPC/Netpackage"
	party "MPC/Party"
	"fmt"
	"math"
	"strconv"
	"sync"
)

type Ceps struct {
	config     *config.Config
	peer       party.IPeer
	shamir     *ShamirSecretSharing
	cMessages  chan netpack.Share
	rShares    rShares
	degree     int
	fShares    fShares
	multInsNum incrementer
}

type incrementer struct {
	nr   int
	lock sync.Mutex
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
	astExp := config.ParseExpression(exp)
	instructionTree, err := config.ConvertAstToTree(astExp)
	if err != nil {
		//TODO maybe shut down peer?
		println(err.Error())
		return 0
	}
	//Send input shares
	toSendIdentifier := netpack.ShareIdentifier{Ins: "p" + strconv.Itoa(int(prot.config.VariableConfig.PartyNr)), PartyNr: int(prot.config.VariableConfig.PartyNr)}
	shares := prot.shamir.makeShares(int64(prot.config.ConstantConfig.NumberOfParties), toSendIdentifier)
	prot.handleShare(shares)

	//Do instructions
	prot.calculateInstruction(*instructionTree)
	finalResult := instructionTree.Instruction.Result

	//output reconstruction
	res := prot.outputReconstruction(finalResult)
	return res
}

func (prot *Ceps) calculateInstruction(instructionTree config.InstructionTree) {
	if instructionTree.Left != nil {
		go prot.calculateInstruction(*instructionTree.Left)
	}
	if instructionTree.Right != nil {
		go prot.calculateInstruction(*instructionTree.Right)
	}
	switch instructionTree.Instruction.Op {
	case config.Add:
		prot.add(*instructionTree.Instruction)
	case config.Multiply:
		prot.multInsNum.lock.Lock()
		prot.multInsNum.nr += 1
		prot.multiply(prot.multInsNum.nr, *instructionTree.Instruction)
		prot.multInsNum.lock.Unlock()
	case config.Scalar:
		prot.scalar(*instructionTree.Instruction)
	}
	return
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

	}
}

func (prot *Ceps) waitForShares(needToWaitOn []netpack.ShareIdentifier) {
	shares := make(map[netpack.ShareIdentifier]*netpack.Share)
	for {
		for _, s := range needToWaitOn {
			prot.rShares.mu.Lock()
			temp := prot.rShares.receivedShares[s]
			prot.rShares.mu.Unlock()
			if temp != nil {
				shares[s] = temp
			}
		}
		if len(shares) == len(needToWaitOn) {
			return
		}
	}

}

func (prot *Ceps) addResultShare(insResult string, value int64) {
	resultShare := &netpack.Share{}
	resultShare.Value = value
	resultShare.Identifier = netpack.ShareIdentifier{Ins: insResult, PartyNr: int(prot.config.VariableConfig.PartyNr)}
	prot.rShares.mu.Lock()
	prot.rShares.receivedShares[resultShare.Identifier] = resultShare
	prot.rShares.mu.Unlock()
}

func (prot *Ceps) handleShare(shares []netpack.Share) {
	fmt.Printf("Party %v send %v\n", prot.config.VariableConfig.PartyNr, shares)
	prot.peer.SendShares(shares)
	prot.rShares.mu.Lock()
	myShare := shares[int(prot.config.VariableConfig.PartyNr)-1]
	prot.rShares.receivedShares[myShare.Identifier] = &myShare
	prot.rShares.mu.Unlock()
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

func (prot *Ceps) add(ins config.Instruction) {
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

func (prot *Ceps) multiply(instructionNumber int, ins config.Instruction) {
	leftIden := prot.createWaitShareIdentifier(ins.Left)
	rightIden := prot.createWaitShareIdentifier(ins.Right)
	prot.waitForShares([]netpack.ShareIdentifier{leftIden, rightIden})
	//This is not the s0 value, but each party's perception of the value that they will use in the new polynomial
	prot.rShares.mu.Lock()
	leftVal := prot.rShares.receivedShares[leftIden].Value
	rightVal := prot.rShares.receivedShares[rightIden].Value
	prot.rShares.mu.Unlock()
	secretValue := prot.shamir.field.Multiply(leftVal, rightVal)
	SSS := makeShamirSecretSharing(secretValue, prot.shamir.field, prot.degree)
	toSendIdentifier := netpack.ShareIdentifier{Ins: "m" + strconv.Itoa(instructionNumber), PartyNr: int(prot.config.VariableConfig.PartyNr)}
	shares := SSS.makeShares(int64(prot.config.ConstantConfig.NumberOfParties), toSendIdentifier)
	prot.handleShare(shares)
	var multiplicationIdentifiers []netpack.ShareIdentifier
	for i := 1; i <= int(prot.config.ConstantConfig.NumberOfParties); i++ {
		multiplicationIdentifiers = append(multiplicationIdentifiers, netpack.ShareIdentifier{Ins: "m" + strconv.Itoa(instructionNumber), PartyNr: i})
	}
	prot.waitForShares(multiplicationIdentifiers)

	var sharesForLagrange []netpack.Share
	for _, i := range multiplicationIdentifiers {
		sharesForLagrange = append(sharesForLagrange, *prot.rShares.receivedShares[i])
	}

	//Use lagrange interpolation on received shares, note that degree needs to be changed to numParties-1
	fmt.Printf("In multiply party %v is calling lagrange with degree: %v, and shares: %v\n", prot.config.VariableConfig.PartyNr, int(prot.config.ConstantConfig.NumberOfParties-1), sharesForLagrange)
	value, _ := SSS.lagrangeInterpolation(sharesForLagrange, int(prot.config.ConstantConfig.NumberOfParties-1))
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) scalar(ins config.Instruction) {
	rightIden := prot.createWaitShareIdentifier(ins.Right)
	prot.waitForShares([]netpack.ShareIdentifier{rightIden})
	scalar, err := strconv.Atoi(ins.Left)
	if err != nil {
		println("Received non-integer as scalar in instruction")
		return
	}
	prot.rShares.mu.Lock()
	varValue := prot.rShares.receivedShares[rightIden].Value
	prot.rShares.mu.Unlock()
	value := prot.shamir.field.Multiply(int64(scalar), varValue)
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
