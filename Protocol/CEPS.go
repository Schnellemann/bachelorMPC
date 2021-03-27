package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	netpack "MPC/Netpackage"
	party "MPC/Party"
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
)

type Ceps struct {
	config    *config.Config
	peer      party.IPeer
	shamir    *ShamirSecretSharing
	cMessages chan netpack.Share
	rShares   rShares
	degree    int
	fShares   fShares
}

type fShares struct {
	finalShares []netpack.Share
	mu          sync.Mutex
}

type rShares struct {
	receivedShares map[string]*netpack.Share
	mu             sync.Mutex
}

func mkProtocol(config *config.Config, secret int64, field field.Field, peer party.IPeer) *Ceps {
	prot := new(Ceps)
	prot.cMessages = make(chan netpack.Share)
	prot.config = config
	prot.peer = peer
	prot.degree = int(math.Ceil(prot.config.ConstantConfig.NumberOfParties/2) - 1)
	prot.shamir = makeShamirSecretSharing(secret, field, prot.degree)
	prot.rShares = rShares{receivedShares: make(map[string]*netpack.Share)}
	return prot
}

func (prot *Ceps) run() int64 {

	partyProgress := make(chan int)
	prot.peer.SetProgress(partyProgress)
	//Start peer
	prot.peer.StartPeer(prot.cMessages)

	//wait group for start peer
	<-partyProgress
	//Start receiving messages from the network
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go prot.receive(ctx)
	//Convert string expression into instruction list
	exp := prot.config.ConstantConfig.Expression
	astExp := config.ParseExpression(exp)
	finalResult, instructionList, err := config.ConvertAstToExpressionList(astExp)
	if err != nil {
		//TODO maybe shut down peer?
		println(err)
		return 0
	}

	//Do instructions
	for i, ins := range instructionList {
		switch ins.Op {
		case config.Add:
			prot.add(ins)
		case config.Multiply:
			prot.multiply(i, ins)
		case config.Scalar:
			prot.scalar(ins)
		}
	}

	//output reconstruction
	res := prot.outputReconstruction(finalResult)

	fmt.Printf("Party %v got %v as the final result\n", int(prot.config.VariableConfig.PartyNr), res)
	return res
}

func (prot *Ceps) receive(ctx context.Context) {
	for {
		select {
		case message := <-prot.cMessages:
			if string(message.Identifier[0]) == "o" {
				prot.fShares.mu.Lock()
				prot.fShares.finalShares = append(prot.fShares.finalShares, message)
				prot.fShares.mu.Unlock()
			} else {
				prot.rShares.mu.Lock()
				prot.rShares.receivedShares[message.Identifier] = &message
				prot.rShares.mu.Unlock()
			}

		case <-ctx.Done():
			fmt.Println("Protocol received shutdown signal, closing messageChannel!")
			close(prot.cMessages)
			fmt.Println("Protocol closed messageChannel")
			return
		}

	}
}

func (prot *Ceps) waitForShares(needToWaitOn []string) {
	shares := make(map[string]*netpack.Share)
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
	resultShare.Identifier = "o" + strconv.Itoa(int(prot.config.VariableConfig.PartyNr))
	prot.rShares.mu.Lock()
	prot.rShares.receivedShares[insResult] = resultShare
	prot.rShares.mu.Unlock()
}

func (prot *Ceps) add(ins config.Instruction) {
	prot.waitForShares([]string{ins.Left, ins.Right})
	prot.rShares.mu.Lock()
	leftVal := prot.rShares.receivedShares[ins.Left].Value
	rightVal := prot.rShares.receivedShares[ins.Right].Value
	prot.rShares.mu.Unlock()
	value := prot.shamir.field.Add(leftVal, rightVal)
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) multiply(instructionNumber int, ins config.Instruction) {
	prot.waitForShares([]string{ins.Left, ins.Right})
	//This is not the s0 value, but each party's perception of the value that they will use in the new polynomial
	prot.rShares.mu.Lock()
	leftVal := prot.rShares.receivedShares[ins.Left].Value
	rightVal := prot.rShares.receivedShares[ins.Right].Value
	prot.rShares.mu.Unlock()
	secretValue := prot.shamir.field.Multiply(leftVal, rightVal)
	SSS := makeShamirSecretSharing(secretValue, prot.shamir.field, prot.degree)
	toSendIdentifier := "m" + strconv.Itoa(instructionNumber) + "," + strconv.Itoa(int(prot.config.VariableConfig.PartyNr))
	shares := SSS.makeShares(int64(prot.config.ConstantConfig.NumberOfParties), toSendIdentifier)

	prot.peer.SendShares(shares)
	var multiplicationIdentifiers []string
	for i := 0; i < int(prot.config.ConstantConfig.NumberOfParties); i++ {
		multiplicationIdentifiers = append(multiplicationIdentifiers, ("m" + strconv.Itoa(instructionNumber) + "," + strconv.Itoa(i)))
	}
	prot.waitForShares(multiplicationIdentifiers)
	//Now we can actually get our value

	//Find the recombination vector

	//Use lagrange interpolation on received shares, note that degree needs to be changed to numParties-1
	value, _ := SSS.lagrangeInterpolation(shares, int(prot.config.ConstantConfig.NumberOfParties-1))
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) scalar(ins config.Instruction) {
	prot.waitForShares([]string{ins.Right})
	scalar, err := strconv.Atoi(ins.Left)
	if err != nil {
		fmt.Printf("Received non-integer as scalar in instruction")
		return
	}
	prot.rShares.mu.Lock()
	varValue := prot.rShares.receivedShares[ins.Right].Value
	prot.rShares.mu.Unlock()
	value := prot.shamir.field.Multiply(int64(scalar), varValue)
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) outputReconstruction(finalResult string) int64 {
	//Send out result share
	prot.rShares.mu.Lock()
	resultShare := prot.rShares.receivedShares[finalResult]
	prot.rShares.mu.Unlock()
	resultShare.Identifier = "o" + strconv.Itoa(int(prot.config.VariableConfig.PartyNr))
	prot.peer.SendFinal(*resultShare)
	prot.fShares.mu.Lock()
	prot.fShares.finalShares = append(prot.fShares.finalShares, *resultShare)
	prot.fShares.mu.Unlock()
	shares := prot.waitForFinalShares()
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
