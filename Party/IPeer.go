package party

import (
	netpack "MPC/Netpackage"
)

type IPeer interface {
	SendShares([]netpack.Share)
	StartPeer(chan netpack.Share)
	SetProgress(chan int)
	SendFinal(netpack.Share)
}
