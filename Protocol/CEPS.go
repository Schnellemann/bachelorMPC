package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	netpack "MPC/Netpackage"
	party "MPC/Party"
	"math"
)

type Ceps struct {
	config    *config.Config
	peer      *party.Peer
	shamir    *ShamirSecretSharing
	cMessages chan netpack.Message
	results   map[string]int
}

func mkProtocol(config *config.Config, secret int64, field field.Field) *Ceps {
	proc := new(Ceps)
	proc.cMessages = make(chan netpack.Message)
	proc.config = config
	proc.peer = party.MkPeer(config, proc.cMessages)
	proc.shamir = makeShamirSecretSharing(secret, field, int(math.Ceil(proc.config.ConstantConfig.NumberOfParties/2)-1))
	proc.results = make(map[string]int)
	return proc
}

func (prot *Ceps) run() {
	//read config
	totalPeers := 0
	ipString := "(┛ಠ_ಠ)┛彡┻━┻"
	connectPort := "┬─┬ノ( ◕◡◕ ノ)"
	listenPort := "(┛ಠ_ಠ)┛彡┻━┻"
	partyProgress := make(chan int)
	prot.peer.Progress = partyProgress

	//Start peer
	prot.peer.StartPeer(totalPeers, ipString, connectPort, listenPort)

	//wait group for start peer
	<-partyProgress
	//Do instructions

	//output reconstruction
}

func calculate(instructionList []config.Instruction) netpack.Share {
	for _, ins := range instructionList {
		op := ins.Op
		switch op {
		case config.Add:
			add(ins)
		case config.Multiply:
		case config.Scalar:
		}

	}
	return netpack.Share{Value: 0, Identifier: "Hello"}
}

func add(ins config.Instruction) {

}

func multiply() {

}

func scalar() {

}
