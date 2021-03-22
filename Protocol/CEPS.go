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
}

func mkProtocol() {
	numberOfParties := 3
	proc := new(Ceps)
	proc.cMessages = make(chan netpack.Message)
	proc.peer = party.MkPeer(numberOfParties, proc.cMessages)
}
