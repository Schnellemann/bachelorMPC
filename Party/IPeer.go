package party

import (
	netpack "MPC/Netpackage"
	"sync"
)

type IPeer interface {
	SendShares([]netpack.Share)
	SendShare(netpack.Share, int)
	StartPeer(chan netpack.Share, *sync.WaitGroup)
	SendFinal(netpack.Share)
}
