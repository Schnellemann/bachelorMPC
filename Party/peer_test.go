package party

import (
	config "MPC/Config"
	netpack "MPC/Netpackage"
	"fmt"
	"reflect"
	"testing"
	"time"
)

var filepath string = "ConfigFiles/configTest.json"

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
	p.StartPeer()
	time.Sleep(3 * time.Second)
	p2.StartPeer()
	time.Sleep(3 * time.Second)
	p3.StartPeer()
	time.Sleep(3 * time.Second)

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
	}
	if len(p3.connections) != 2 {
		t.Errorf(assertEqualError(len(p3.cConnections), 2))
	}

}

func contains(s []netpack.PeerTuple, e netpack.PeerTuple) bool {
	for _, p := range s {
		if p == e {
			return true
		}
	}
	return false
}

func TestConnectionlist(t *testing.T) {
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

	fmt.Printf("The connection %v holds: \n", p.Number)
	for _, con := range p.connections {
		fmt.Printf("%v \n", con.Number)
	}

	fmt.Printf("The connection %v holds: \n", p2.Number)
	for _, con := range p2.connections {
		fmt.Printf("%v \n", con.Number)
	}

	fmt.Printf("The connection %v holds: \n", p3.Number)
	for _, con := range p3.connections {
		fmt.Printf("%v \n", con.Number)
	}
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

	shouldHold := []netpack.PeerTuple{{"127.0.1.1:40002", 1}, {"127.0.1.1:6970", 2}, {"127.0.1.1:6971", 3}}

	for i := 0; i < 3; i++ {
		for _, j := range shouldHold {
			if !contains(peers[i].peerlist.ipPorts, j) {
				t.Errorf("peer %v does not hold peer %v's PeerTuple", i+1, j.Number)
			}
		}
	}
}
