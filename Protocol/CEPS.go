package protocol

import (
	party "MPC/party"
	netpack "MPC/netPackage"
	config "MPC/config"
)

type Ceps struct {
	config	Config
	peer   party.Peer
	shamir ShamirSecretSharing
	cMessages	Message
}

func mkProtocol() {

	proc := new(Ceps)
	proc.cMessages = make(chan Message)
	proc.peer = party.mkPeer(numberOfParties, proc.cMessages)
}
