package party

import (
	netpack "MPC/Netpackage"
	"sync"
)

type IPeer interface {
	SendShares([]netpack.Share)
	StartPeer(chan netpack.Share, *sync.WaitGroup)
	SendFinal(netpack.Share)
}
