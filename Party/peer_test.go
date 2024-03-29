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
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Started peer 2")
	p2.StartPeer(nil, &wg)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Started peer 3")
	p3.StartPeer(nil, &wg)
	time.Sleep(100 * time.Millisecond)
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
		time.Sleep(100 * time.Millisecond)
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
	configs := config.MakeConfigs("127.0.1.1", "", []int{1, 2, 3})
	conf := configs[0]
	conf2 := configs[1]
	conf3 := configs[2]
	p := MkPeer(conf)
	p2 := MkPeer(conf2)
	p3 := MkPeer(conf3)
	pChan1 := make(chan netpackage.Share)
	pChan2 := make(chan netpackage.Share)
	pChan3 := make(chan netpackage.Share)
	/*
		Connect them
	*/
	var wg sync.WaitGroup
	wg.Add(3)
	p.StartPeer(pChan1, &wg)
	time.Sleep(100 * time.Millisecond)
	p2.StartPeer(pChan2, &wg)
	time.Sleep(100 * time.Millisecond)
	p3.StartPeer(pChan3, &wg)
	wg.Wait()
	peers := []*Peer{p, p2, p3}
	shouldHold := []string{"127.0.1.1:50000", "127.0.1.1:50010", "127.0.1.1:50020"}

	for i := 0; i < 3; i++ {
		for _, j := range shouldHold {
			if !contains(peers[i].peerlist.ipPorts, j) {
				t.Errorf("peer %v does not hold peer %v ip", i+1, j)
			}
		}
	}
}
