package party

import (
	config "MPC/Config"
	netpackage "MPC/Netpackage"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBandwidthConstraint(t *testing.T) {
	configs := config.ReadConfig(filepath)
	conf := configs[0]
	conf2 := configs[1]
	conf3 := configs[2]
	conf4 := configs[3]
	conf5 := configs[4]

	peer := MkPeer(conf)
	peer2 := MkPeer(conf2)
	peer3 := MkPeer(conf3)
	peer4 := MkPeer(conf4)
	peer5 := MkPeer(conf5)

	p := MkBandwidthPeer(conf, peer, 10, 10*time.Millisecond)
	p2 := MkBandwidthPeer(conf2, peer2, 10, 10*time.Millisecond)
	p3 := MkBandwidthPeer(conf3, peer3, 10, 10*time.Millisecond)
	p4 := MkBandwidthPeer(conf4, peer4, 10, 10*time.Millisecond)
	p5 := MkBandwidthPeer(conf5, peer5, 10, 10*time.Millisecond)
	/*
		Make channels for message
	*/
	pChan1 := make(chan netpackage.Share)
	pChan2 := make(chan netpackage.Share)
	pChan3 := make(chan netpackage.Share)
	pChan4 := make(chan netpackage.Share)
	pChan5 := make(chan netpackage.Share)

	/*
		Connect them
	*/
	var wg sync.WaitGroup
	wg.Add(5)
	fmt.Println("Started peer 1")
	p.StartPeer(pChan1, &wg)
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Started peer 2")
	p2.StartPeer(pChan2, &wg)
	time.Sleep(1 * time.Second)
	fmt.Println("Started peer 3")
	p3.StartPeer(pChan3, &wg)
	time.Sleep(1 * time.Second)
	fmt.Println("Started peer 4")
	p4.StartPeer(pChan4, &wg)
	time.Sleep(1 * time.Second)
	fmt.Println("Started peer 5")
	p5.StartPeer(pChan5, &wg)
	time.Sleep(1 * time.Second)
	wg.Wait()
	shares := []netpackage.Share{
		{Value: 1, Identifier: netpackage.ShareIdentifier{Ins: "share1", PartyNr: 1}},
		{Value: 2, Identifier: netpackage.ShareIdentifier{Ins: "share2", PartyNr: 1}},
		{Value: 3, Identifier: netpackage.ShareIdentifier{Ins: "share3", PartyNr: 1}},
		{Value: 4, Identifier: netpackage.ShareIdentifier{Ins: "share4", PartyNr: 1}},
		{Value: 5, Identifier: netpackage.ShareIdentifier{Ins: "share5", PartyNr: 1}},
	}
	startSend := time.Now()
	fmt.Println("Started sending")
	go p.SendShares(shares)
	<-pChan1
	<-pChan2
	<-pChan3
	<-pChan4
	<-pChan5
	endSend := time.Now()
	timeToSend := endSend.Sub(startSend)
	if (timeToSend) < 190*time.Millisecond {
		t.Errorf("The sending of the message was too fast, got %v", timeToSend)
	}
}
