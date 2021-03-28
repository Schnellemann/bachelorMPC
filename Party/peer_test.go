package party

import (
	config "MPC/Config"
	netpackage "MPC/Netpackage"
	"fmt"
	"reflect"
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
	constantConfig := config.ConstantConfig{"", 3, []string{"127.0.1.1:40002", "127.0.1.1:60716", "127.0.1.1:60817"}}
	variableConfig1 := config.VariableConfig{"127.0.1.1:40002", "", 1, 1}
	variableConfig2 := config.VariableConfig{"127.0.1.1:60716", "127.0.1.1:40002", 2, 2}
	variableConfig3 := config.VariableConfig{"127.0.1.1:60817", "127.0.1.1:60716", 3, 3}
	/*
		Make the peers
	*/
	configs := []config.Config{{VariableConfig: variableConfig1, ConstantConfig: constantConfig},
		{VariableConfig: variableConfig2, ConstantConfig: constantConfig},
		{VariableConfig: variableConfig3, ConstantConfig: constantConfig}}
	conf := &configs[0]
	conf2 := &configs[1]
	conf3 := &configs[2]
	p := MkPeer(conf)
	p2 := MkPeer(conf2)
	p3 := MkPeer(conf3)

	/*
		Connect them
	*/
	fmt.Println("Started peer 1")
	p.StartPeer(nil)
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Started peer 2")
	p2.StartPeer(nil)
	time.Sleep(1 * time.Second)
	fmt.Println("Started peer 3")
	p3.StartPeer(nil)
	time.Sleep(1 * time.Second)

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
	if len(p.connections) != 2 {
		t.Errorf(assertEqualError(len(p.connections), 2))
	}
	if len(p2.connections) != 2 {
		t.Errorf(assertEqualError(len(p2.connections), 2))
	}
	if len(p3.connections) != 2 {
		t.Errorf(assertEqualError(len(p3.connections), 2))
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
	progress1 := make(chan int)
	progress2 := make(chan int)
	progress3 := make(chan int)
	progress4 := make(chan int)
	progress5 := make(chan int)
	p.SetProgress(progress1)
	p2.SetProgress(progress2)
	p3.SetProgress(progress3)
	p4.SetProgress(progress4)
	p5.SetProgress(progress5)
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
	fmt.Println("Started peer 1")
	p.StartPeer(pChan1)
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Started peer 2")
	p2.StartPeer(pChan2)
	time.Sleep(1 * time.Second)
	fmt.Println("Started peer 3")
	p3.StartPeer(pChan3)
	time.Sleep(1 * time.Second)
	fmt.Println("Started peer 4")
	p4.StartPeer(pChan4)
	time.Sleep(1 * time.Second)
	fmt.Println("Started peer 5")
	p5.StartPeer(pChan5)
	time.Sleep(1 * time.Second)
	<-progress1
	<-progress2
	<-progress3
	<-progress4
	<-progress5
	shares := []netpackage.Share{{1, "share1"}, {2, "share2"}, {3, "share3"}, {4, "share4"}, {5, "share5"}}

	p.SendShares(shares)
	share2Res := <-pChan2
	share3Res := <-pChan3
	share4Res := <-pChan4
	share5Res := <-pChan5

	if share2Res.Value != 2 && share2Res.Identifier != "share2" {
		t.Errorf("Wrong share recieved at peer2, should have value: %v, got: %v", 2, share2Res.Value)
	}
	if share3Res.Value != 3 && share3Res.Identifier != "share3" {
		t.Errorf("Wrong share recieved at peer3, should have value: %v, got: %v", 3, share3Res.Value)
	}
	if share3Res.Value != 4 && share4Res.Identifier != "share4" {
		t.Errorf("Wrong share recieved at peer3, should have value: %v, got: %v", 4, share4Res.Value)
	}
	if share3Res.Value != 5 && share5Res.Identifier != "share5" {
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
	p.StartPeer(nil)
	time.Sleep(3 * time.Second)
	p2.StartPeer(nil)
	time.Sleep(3 * time.Second)
	p3.StartPeer(nil)
	time.Sleep(3 * time.Second)

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
