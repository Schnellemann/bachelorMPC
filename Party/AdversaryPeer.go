package party

import (
	netpack "MPC/Netpackage"
	"sync"
)

type AdversaryPeer struct {
	peer              Peer
	toProtocolChannel chan netpack.Share
}

func (p *AdversaryPeer) SendShares(shares []netpack.Share) {
	//Change such that it gets all shares (including its own)
	//Send to collector
	p.peer.SendShares(shares)

}

func (p *AdversaryPeer) StartPeer(messagechannel chan netpack.Share, wg *sync.WaitGroup) {
	p.toProtocolChannel = messagechannel
	peerChannel := make(chan netpack.Share)
	p.peer.StartPeer(peerChannel, wg)

}

func (p *AdversaryPeer) SendFinal(share netpack.Share) {
	p.peer.SendFinal(share)
	//Send to collector
}

func (p *AdversaryPeer) receive(peerChannel chan netpack.Share) {
	for {
		message := <-peerChannel
		p.toProtocolChannel <- message
		//TODO handle adversary sending
		//Send to collector
	}
}
