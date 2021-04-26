package party

import (
	config "MPC/Config"
	netpack "MPC/Netpackage"
	"sync"
	"time"
)

type DelayedPeer struct {
	Peer  IPeer
	delay time.Duration
}

func MkDelayedPeer(config *config.Config, delay time.Duration, peer IPeer) *DelayedPeer {
	cp := new(DelayedPeer)
	cp.Peer = peer
	cp.delay = delay
	return cp
}

func (cnP *DelayedPeer) SendShares(shares []netpack.Share) {
	time.Sleep(cnP.delay)
	cnP.Peer.SendShares(shares)
}

func (cnP *DelayedPeer) SendShare(share netpack.Share, nr int) {
	time.Sleep(cnP.delay)
	cnP.Peer.SendShare(share, nr)
}

func (cnP *DelayedPeer) StartPeer(shareChannel chan netpack.Share, waitGroup *sync.WaitGroup) {
	cnP.Peer.StartPeer(shareChannel, waitGroup)
}

func (cnP *DelayedPeer) SendFinal(finalShare netpack.Share) {
	time.Sleep(cnP.delay)
	cnP.Peer.SendFinal(finalShare)
}
