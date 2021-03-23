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
}

func mkProtocol(config *config.Config, secret int64, field field.Field) *Ceps {
	proc := new(Ceps)
	proc.cMessages = make(chan netpack.Message)
	proc.config = config
	proc.peer = party.MkPeer(config, proc.cMessages)
	proc.shamir = makeShamirSecretSharing(secret, field, int(math.Ceil(proc.config.NumberOfParties/2)-1))
	return proc
}
