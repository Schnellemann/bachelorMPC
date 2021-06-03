package party

import (
	netpack "MPC/Netpackage"
	"sync"
)

/*
	Network interface to use in a secret sharings scheme
*/
type IPeer interface {
	SendShares([]netpack.Share)
	SendShare(netpack.Share, int)
	StartPeer(chan netpack.Share, *sync.WaitGroup)
	SendFinal(netpack.Share)
}
