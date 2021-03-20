package party

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var ip string = "192.168.0.147"

func assertEqualError(received interface{}, expected interface{}) string {
	return fmt.Sprintf("Received %v (type %v), expected %v (type %v)", received, reflect.TypeOf(received), expected, reflect.TypeOf(expected))
}

func TestConnection(t *testing.T) {
	/*
		Make the peers
	*/
	p := mkPeer(1)
	p2 := mkPeer(2)
	p3 := mkPeer(3)

	/*
		Connect them
	*/
	p.startPeer(3, ip, "", "61515")
	time.Sleep(1000 * time.Millisecond)
	p2.startPeer(3, ip, "61515", "60516")
	time.Sleep(1000 * time.Millisecond)
	p3.startPeer(3, ip, "61515", "60417")
	time.Sleep(1000 * time.Millisecond)

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

func contains(s []PeerTuple, e PeerTuple) bool {
	for _, p := range s {
		if p == e {
			return true
		}
	}
	return false
}

func TestPeerlists(t *testing.T) {
	p := mkPeer(1)
	p2 := mkPeer(2)
	p3 := mkPeer(3)
	peers := []Peer{*p, *p2, *p3}

	p.startPeer(3, ip, "", "61515")
	time.Sleep(1000 * time.Millisecond)
	p2.startPeer(3, ip, "61515", "60516")
	time.Sleep(1000 * time.Millisecond)
	p3.startPeer(3, ip, "61515", "60417")
	time.Sleep(1000 * time.Millisecond)

	shouldHold := []PeerTuple{{ip + ":" + "61515", 1}, {ip + ":" + "60516", 2}, {ip + ":" + "60417", 3}}

	for i := 0; i < 3; i++ {
		for _, j := range shouldHold {
			if !contains(peers[i].peerlist.ipPorts, j) {
				t.Errorf("peer %v does not hold peer %v's PeerTuple", i+1, j.Number)
			}
		}
	}
}
