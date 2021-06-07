package party

import (
	config "MPC/Config"
	netpack "MPC/Netpackage"
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"sync"
	"time"
)

type BandwidthPeer struct {
	Peer           IPeer
	SendChannel    chan netpack.Share
	RecieveChannel chan netpack.Share
	width          int
	penalty        time.Duration
}

func MkBandwidthPeer(config *config.Config, Peer IPeer, width int, penalty time.Duration) *BandwidthPeer {
	bp := new(BandwidthPeer)
	bp.RecieveChannel = make(chan netpack.Share)
	bp.penalty = penalty
	bp.width = width
	bp.Peer = Peer
	return bp
}

func (bp *BandwidthPeer) SendShares(shares []netpack.Share) {
	bp.Peer.SendShares(shares)
}

func (bp *BandwidthPeer) SendShare(share netpack.Share, nr int) {
	bp.Peer.SendShare(share, nr)
}

func (bp *BandwidthPeer) StartPeer(shareChannel chan netpack.Share, waitGroup *sync.WaitGroup) {
	bp.SendChannel = shareChannel
	go bp.propagateMessage()
	bp.Peer.StartPeer(bp.RecieveChannel, waitGroup)
}

func (bp *BandwidthPeer) SendFinal(finalShare netpack.Share) {
	bp.Peer.SendFinal(finalShare)
}

func (bp *BandwidthPeer) propagateMessage() {
	for {
		message := <-bp.RecieveChannel
		newPeerPackage := new(netpack.NetPackage)
		newPeerPackage.Share = message
		bytearray, err := bp.marshalBinary(newPeerPackage)
		if err != nil {
			fmt.Println("BandwidthPeer: Message to byte conversion fail")
		}
		messageSize := len(bytearray)
		if messageSize > bp.width {
			toWait := int64(math.Ceil((float64(messageSize) / float64(bp.width)) - 1))
			time.Sleep(time.Duration(toWait) * bp.penalty)
		}
		bp.SendChannel <- message
	}
}

func (bp *BandwidthPeer) marshalBinary(pack *netpack.NetPackage) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(pack); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
