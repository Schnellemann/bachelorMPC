package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	netpack "MPC/Netpackage"
	party "MPC/Party"
	"context"
	"fmt"
	"math"
)

type Ceps struct {
	config              *config.Config
	peer                *party.Peer
	shamir              *ShamirSecretSharing
	cMessages           chan netpack.Message
	intermediaryResults map[int][]netpack.Share
}

func mkProtocol(config *config.Config, secret int64, field field.Field) *Ceps {
	proc := new(Ceps)
	proc.cMessages = make(chan netpack.Message)
	proc.config = config
	proc.peer = party.MkPeer(config, proc.cMessages)
	proc.shamir = makeShamirSecretSharing(secret, field, int(math.Ceil(proc.config.ConstantConfig.NumberOfParties/2)-1))
	proc.intermediaryResults = make(map[int][]netpack.Share)
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

func (prot *Ceps) add(ins config.Instruction) {

}

func (prot *Ceps) multiply(instructionNumber int, ins config.Instruction) {

}

func (prot *Ceps) scalar(ins config.Instruction) {

}

func (prot *Ceps) outputReconstruction(finalResultName string) int {
	return 0
}
