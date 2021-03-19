package party

import (
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	/*
		Make the peers
	*/
	p := mkPeer()
	p2 := mkPeer()
	p3 := mkPeer()

	/*
		Connect them
	*/
	p.startPeer(ip, "", "61515")
	time.Sleep(1000 * time.Millisecond)
	p2.startPeer(ip, "61515", "60516")
	time.Sleep(1000 * time.Millisecond)
	p3.startPeer(ip, "61515", "60417")
	time.Sleep(1000 * time.Millisecond)

	AssertEqual(t, len(p.localPeers.ipPorts), 3)
	AssertEqual(t, len(p2.localPeers.ipPorts), 3)
	AssertEqual(t, len(p3.localPeers.ipPorts), 3)
	AssertEqual(t, len(p.connections), 2)
	AssertEqual(t, len(p2.connections), 2)
	AssertEqual(t, len(p3.connections), 2)

}
