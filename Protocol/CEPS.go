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
)

type Ceps struct {
	config         *config.Config
	peer           *party.Peer
	shamir         *ShamirSecretSharing
	cMessages      chan netpack.Message
	receivedShares map[string]*netpack.Share
}

func mkProtocol(config *config.Config, secret int64, field field.Field) *Ceps {
	proc := new(Ceps)
	proc.cMessages = make(chan netpack.Message)
	proc.config = config
	proc.peer = party.MkPeer(config, proc.cMessages)
	proc.shamir = makeShamirSecretSharing(secret, field, int(math.Ceil(proc.config.ConstantConfig.NumberOfParties/2)-1))
	proc.receivedShares = make(map[string]*netpack.Share)
	return proc
}

func (prot *Ceps) run() int {

	partyProgress := make(chan int)
	prot.peer.Progress = partyProgress
	//Start peer
	prot.peer.StartPeer()

	//wait group for start peer
	<-partyProgress
	//Start receiving messages from the network
	ctx, cancelFunc := context.WithCancel(context.Background())

	go prot.receive(ctx)
	//Convert string expression into instruction list
	exp := prot.config.ConstantConfig.Expression
	astExp := config.ParseExpression(exp)
	finalResultName, instructionList, err := config.ConvertAstToExpressionList(astExp)
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
	res := prot.outputReconstruction(finalResultName)
	cancelFunc()
	return res
}

func (prot *Ceps) receive(ctx context.Context) {
	for {
		select {
		case message := <-prot.cMessages:
			//TODO

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
			temp := prot.receivedShares[s]
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
	prot.receivedShares[insResult] = resultShare
}

func (prot *Ceps) add(ins config.Instruction) {
	prot.waitForShares([]string{ins.Left, ins.Right})
	value := prot.shamir.field.Add(prot.receivedShares[ins.Left].Value, prot.receivedShares[ins.Right].Value)
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) multiply(instructionNumber int, ins config.Instruction) {
	prot.waitForShares([]string{ins.Left, ins.Right})
	//This is not the s0 value, but each party's perception of the value that they will use in the new polynomial
	secretValue := prot.shamir.field.Multiply(prot.receivedShares[ins.Left].Value, prot.receivedShares[ins.Right].Value)
	SSS := makeShamirSecretSharing(secretValue, prot.shamir.field, prot.shamir.degree)
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
	//value := lagrangeInterpolation()
	value := int64(0)
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) scalar(ins config.Instruction) {
	prot.waitForShares([]string{ins.Right})
	scalar, err := strconv.Atoi(ins.Left)
	if err != nil {
		fmt.Printf("Received non-integer as scalar in instruction")
		return
	}
	value := prot.shamir.field.Multiply(int64(scalar), prot.receivedShares[ins.Right].Value)
	prot.addResultShare(ins.Result, value)
}

func (prot *Ceps) outputReconstruction(finalResultName string) int {
	return 0
}
