package party

import (
	config "MPC/Config"
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
		Make the peers
	*/
	configs := config.ReadConfig(filepath)
	conf := &configs[0]
	conf2 := &configs[1]
	conf3 := &configs[2]
	p := MkPeer(conf, nil)
	p2 := MkPeer(conf2, nil)
	p3 := MkPeer(conf3, nil)

	/*
		Connect them
	*/
	fmt.Println("Started peer 1")
	p.StartPeer()
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Started peer 2")
	p2.StartPeer()
	time.Sleep(2 * time.Second)
	fmt.Println("Started peer 3")
	p3.StartPeer()
	time.Sleep(10 * time.Second)

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

		t.Errorf(assertEqualError(len(p.cConnections), 2))
	}
	if len(p2.connections) != 2 {
		t.Errorf(assertEqualError(len(p2.cConnections), 2))
		fmt.Printf("P2 connections is: %v", p.connections)
	}
	if len(p3.connections) != 2 {
		t.Errorf(assertEqualError(len(p3.cConnections), 2))
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
	p := MkPeer(conf, nil)
	p2 := MkPeer(conf2, nil)
	p3 := MkPeer(conf3, nil)

	/*
		Connect them
	*/
	p.StartPeer()
	time.Sleep(3 * time.Second)
	p2.StartPeer()
	time.Sleep(3 * time.Second)
	p3.StartPeer()
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
