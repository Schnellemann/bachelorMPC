package party

import (
	config "MPC/Config"
	netpackage "MPC/Netpackage"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

var filepath string = "peerTestConfig.json"

func assertEqualError(received interface{}, expected interface{}) string {
	return fmt.Sprintf("Received %v (type %v), expected %v (type %v)", received, reflect.TypeOf(received), expected, reflect.TypeOf(expected))
}

func TestConnections(t *testing.T) {

	/*
		Configs
	*/
	constantConfig := config.ConstantConfig{Expression: "", NumberOfParties: 3, Ipports: []string{"127.0.1.1:40002", "127.0.1.1:60716", "127.0.1.1:60817"}}
	variableConfig1 := config.VariableConfig{ListenIpPort: "127.0.1.1:40002", ConnectIpPort: "", PartyNr: 1, Secret: 1}
	variableConfig2 := config.VariableConfig{ListenIpPort: "127.0.1.1:60716", ConnectIpPort: "127.0.1.1:40002", PartyNr: 2, Secret: 2}
	variableConfig3 := config.VariableConfig{ListenIpPort: "127.0.1.1:60817", ConnectIpPort: "127.0.1.1:60716", PartyNr: 3, Secret: 3}
	configs := []config.Config{{VariableConfig: variableConfig1, ConstantConfig: constantConfig},
		{VariableConfig: variableConfig2, ConstantConfig: constantConfig},
		{VariableConfig: variableConfig3, ConstantConfig: constantConfig}}
	/*
		Make the peers
	*/

	conf := &configs[0]
	conf2 := &configs[1]
	conf3 := &configs[2]
	p := MkPeer(conf)
	p2 := MkPeer(conf2)
	p3 := MkPeer(conf3)

	/*
		Connect them
	*/
	var wg sync.WaitGroup
	wg.Add(3)
	fmt.Println("Started peer 1")
	p.StartPeer(nil, &wg)
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Started peer 2")
	p2.StartPeer(nil, &wg)
	time.Sleep(1 * time.Second)
	fmt.Println("Started peer 3")
	p3.StartPeer(nil, &wg)
	time.Sleep(1 * time.Second)
	wg.Wait()
	if len(p.peerlist.ipPorts) != 3 {
		t.Errorf(assertEqualError(len(p.peerlist.ipPorts), 3))
		fmt.Printf("peerlist for p1: %v", p.peerlist.ipPorts)
	}
	if len(p2.peerlist.ipPorts) != 3 {
		t.Errorf(assertEqualError(len(p2.peerlist.ipPorts), 3))
		fmt.Printf("peerlist for p2: %v", p2.peerlist.ipPorts)
	}
	if len(p3.peerlist.ipPorts) != 3 {
		t.Errorf(assertEqualError(len(p3.peerlist.ipPorts), 3))
		fmt.Printf("peerlist for p3: %v", p3.peerlist.ipPorts)
	}
	if len(p.connections.c) != 2 {
		t.Errorf(assertEqualError(len(p.connections.c), 2))
	}
	if len(p2.connections.c) != 2 {
		t.Errorf(assertEqualError(len(p2.connections.c), 2))
	}
	if len(p3.connections.c) != 2 {
		t.Errorf(assertEqualError(len(p3.connections.c), 2))
	}

}

func TestManyConnections(t *testing.T) {

	/*
		Configs
	*/
	configs := config.MakeConfigs("127.0.1.1", "", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	/*
		Make the peers
	*/
	var peers []*Peer
	for _, c := range configs {
		peers = append(peers, MkPeer(c))
	}

	/*
		Connect them
	*/
	var wg sync.WaitGroup
	wg.Add(10)
	for i, p := range peers {
		fmt.Println("Started peer " + strconv.Itoa(i+1))
		p.StartPeer(nil, &wg)
		time.Sleep(200 * time.Millisecond)
	}
	wg.Wait()
	for i, p := range peers {
		if len(p.peerlist.ipPorts) != 10 {
			t.Errorf(assertEqualError(len(p.peerlist.ipPorts), 10))
			fmt.Printf("peerlist for p%v: %v", i+1, p.peerlist.ipPorts)
		}

		if len(p.connections.c) != 9 {
			t.Errorf(assertEqualError(len(p.connections.c), 9))
		}
	}
}

func TestSendShares(t *testing.T) {
	configs := config.ReadConfig(filepath)
	conf := &configs[0]
	conf2 := &configs[1]
	conf3 := &configs[2]
	conf4 := &configs[3]
	conf5 := &configs[4]
	p := MkPeer(conf)
	p2 := MkPeer(conf2)
	p3 := MkPeer(conf3)
	p4 := MkPeer(conf4)
	p5 := MkPeer(conf5)
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
	fmt.Println("Still waiting")
	wg.Wait()
	fmt.Println("Done waiting")
	shares := []netpackage.Share{
		{Value: 1, Identifier: netpackage.ShareIdentifier{Ins: "share1", PartyNr: 1}},
		{Value: 2, Identifier: netpackage.ShareIdentifier{Ins: "share2", PartyNr: 1}},
		{Value: 3, Identifier: netpackage.ShareIdentifier{Ins: "share3", PartyNr: 1}},
		{Value: 4, Identifier: netpackage.ShareIdentifier{Ins: "share4", PartyNr: 1}},
		{Value: 5, Identifier: netpackage.ShareIdentifier{Ins: "share5", PartyNr: 1}},
	}

	go p.SendShares(shares)
	share1Res := <-pChan1
	share2Res := <-pChan2
	share3Res := <-pChan3
	share4Res := <-pChan4
	share5Res := <-pChan5

	if share1Res.Value != 1 && share1Res.Identifier.Ins != "share1" {
		t.Errorf("Wrong share recieved at peer2, should have value: %v, got: %v", 1, share2Res.Value)
	}
	if share2Res.Value != 2 && share2Res.Identifier.Ins != "share2" {
		t.Errorf("Wrong share recieved at peer2, should have value: %v, got: %v", 2, share2Res.Value)
	}
	if share3Res.Value != 3 && share3Res.Identifier.Ins != "share3" {
		t.Errorf("Wrong share recieved at peer3, should have value: %v, got: %v", 3, share3Res.Value)
	}
	if share3Res.Value != 4 && share4Res.Identifier.Ins != "share4" {
		t.Errorf("Wrong share recieved at peer3, should have value: %v, got: %v", 4, share4Res.Value)
	}
	if share3Res.Value != 5 && share5Res.Identifier.Ins != "share5" {
		t.Errorf("Wrong share recieved at peer3, should have value: %v, got: %v", 5, share5Res.Value)
	}
}

func contains(s []string, e string) bool {
	for _, p := range s {
		if p == e {
			return true
		}
	}
	return false
}

func TestPeerlists(t *testing.T) {
	/*
		Make the peers
	*/
	configs := config.ReadConfig(filepath)
	conf := &configs[0]
	conf2 := &configs[1]
	conf3 := &configs[2]
	p := MkPeer(conf)
	p2 := MkPeer(conf2)
	p3 := MkPeer(conf3)

	/*
		Connect them
	*/
	var wg sync.WaitGroup
	wg.Add(3)
	p.StartPeer(nil, &wg)
	time.Sleep(3 * time.Second)
	p2.StartPeer(nil, &wg)
	time.Sleep(3 * time.Second)
	p3.StartPeer(nil, &wg)
	time.Sleep(3 * time.Second)
	wg.Wait()
	peers := []Peer{*p, *p2, *p3}
	shouldHold := []string{"127.0.1.1:40002", "127.0.1.1:6970", "127.0.1.1:6971"}

	for i := 0; i < 3; i++ {
		for _, j := range shouldHold {
			if !contains(peers[i].peerlist.ipPorts, j) {
				t.Errorf("peer %v does not hold peer %v ip", i+1, j)
			}
		}
	}
}
