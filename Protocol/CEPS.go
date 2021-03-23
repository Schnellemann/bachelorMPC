package protocol

import (
	config "MPC/Config"
	netpack "MPC/Netpackage"
	party "MPC/Party"
)

type Ceps struct {
	config    config.Config
	peer      *party.Peer
	shamir    ShamirSecretSharing
	cMessages chan netpack.Message
	results   map[string]int
}

func mkProtocol(numberOfParties int) *Ceps {
	proc := new(Ceps)
	proc.cMessages = make(chan netpack.Message)
	proc.peer = party.MkPeer(numberOfParties, proc.cMessages)
	proc.results = make(map[string]int)
	return proc
}

func (prot *Ceps) run() {
	//read config
	number := 0
	//Start peer (go routine)
	p := party.MkPeer(number, prot.cMessages)
	totalPeers := 0
	ipString := "(┛ಠ_ಠ)┛彡┻━┻"
	connectPort := "┬─┬ノ( ◕◡◕ ノ)"
	listenPort := "(┛ಠ_ಠ)┛彡┻━┻"
	partyProgress := make(chan int)
	p.Progress = partyProgress
	p.StartPeer(totalPeers, ipString, connectPort, listenPort)

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
