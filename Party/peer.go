package party

import (
	aux "MPC/Auxiliary"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
)

type Peer struct {
	Number       int
	ipListen     string
	cConnections chan net.Conn
	cPackages    chan *NetPackage
	connections  []*gob.Encoder
	peerlist     *peerList
}

type peerList struct {
	ipPorts []PeerTuple
	lock    sync.Mutex
}

type PeerTuple struct {
	IpPort string
	Number int
}

func mkPeerList() *peerList {
	pl := new(peerList)
	return pl
}

func mkPeer(number int) *Peer {
	p := new(Peer)
	p.Number = number
	p.peerlist = mkPeerList()
	p.cPackages = make(chan *NetPackage)
	p.cConnections = make(chan net.Conn)

	return p
}

func (p *Peer) handleConnection(dec *gob.Decoder) {
	//A new peer has connected to us
	//Start receiving packages
	for {
		receivedPackage := &NetPackage{}
		err := dec.Decode(receivedPackage)
		if err != nil {
			fmt.Println("Could not decode package from peer")
			fmt.Println(err)
			continue
		} else {
			//If we receive IpPorts we should ignore it, we handle this
			//Only if we're actually waiting for the peerlist
			if receivedPackage.IpPorts == nil {
				p.cPackages <- receivedPackage
			}
		}
	}
}

func (p *Peer) receiveFromChannels() {
	for {
		select {
		case newConnection := <-p.cConnections:
			encoder := gob.NewEncoder(newConnection)
			decoder := gob.NewDecoder(newConnection)
			p.connections = append(p.connections, encoder)
			//send peers to the new connections
			p.sendPeerlist(encoder)
			go p.handleConnection(decoder)

		case newPackage := <-p.cPackages:
			p.processPackage(newPackage)
		}
	}
}

func (p *Peer) sendPeerlist(encoder *gob.Encoder) {
	peerPackage := new(NetPackage)
	peerPackage.IpPorts = p.peerlist.ipPorts
	p.write(encoder, peerPackage)
}

func (p *Peer) write(encoder *gob.Encoder, pack *NetPackage) {
	err := encoder.Encode(pack)
	if err != nil {
		fmt.Println("Could not encode transaction trying again")
		fmt.Println(err)
		p.write(encoder, pack)
	}
}

func (p *Peer) processPackage(pack *NetPackage) {
	if pack.Peer.IpPort != "" {
		//New peer wants to connect to us
		p.peerlist.lock.Lock()
		defer p.peerlist.lock.Unlock()
		p.peerlist.ipPorts = append(p.peerlist.ipPorts, pack.Peer)
	} else {
		//Message
		m := pack.Message
		p.processMessage(m)
	}

}

func (p *Peer) processMessage(message Message) {

}

func (p *Peer) connectToPeers(initialConn string) {
	for _, peer := range p.peerlist.ipPorts {
		//Make sure you don't connect to the initial peer again
		if peer.IpPort != initialConn && peer.IpPort != p.ipListen {
			conn, err := net.Dial("tcp", peer.IpPort)
			if err != nil {
				fmt.Println("Failed to connect to " + peer.IpPort)
				fmt.Println(err)
			} else {
				p.cConnections <- conn
			}
		}
	}
}

func (p *Peer) receivePeers(dec *gob.Decoder) {
	recievedPeersPackage := NetPackage{}
	err := dec.Decode(&recievedPeersPackage)
	if err != nil {
		fmt.Println("Could not decode peer list package from peer")
		return
	}
	p.peerlist.lock.Lock()
	defer p.peerlist.lock.Unlock()
	p.peerlist.ipPorts = append(p.peerlist.ipPorts, recievedPeersPackage.IpPorts...)
}

func (p *Peer) listenForConnections(totalPeers int, ip string, listenPort string) {
	var ipport = ip + ":" + listenPort
	li, err := net.Listen("tcp", ipport)
	if err != nil {
		fmt.Println("Error in listening")
		return
	}
	defer li.Close()
	name, _ := os.Hostname()
	_, port, _ := net.SplitHostPort(li.Addr().String())
	addrs, _ := net.LookupHost(name)
	fmt.Println("Other peers can connect to me on the following ip:port")
	for _, addr := range addrs {
		if aux.IsIpv4Regex(addr) {
			fmt.Println("Address " + ": " + addr + ":" + port)
			p.broadcastPeer(addr + ":" + port)
		}
	}
	i := 1
	for i < totalPeers {
		conn, err := li.Accept()
		if err != nil {
			fmt.Println("Failed connection on accept")
			return
		}
		p.cConnections <- conn
		i++
	}
}

func (p *Peer) broadcastPeer(ipPort string) {
	newPeerPackage := new(NetPackage)
	newPeerPackage.Peer = PeerTuple{IpPort: ipPort, Number: p.Number}
	for _, c := range p.connections {
		p.write(c, newPeerPackage)
	}
}

func (p *Peer) startPeer(totalPeers int, ip string, connectToPort string, listenOnPort string) {
	p.ipListen = ip + ":" + listenOnPort
	p.peerlist.lock.Lock()
	p.peerlist.ipPorts = append(p.peerlist.ipPorts, PeerTuple{ip + ":" + listenOnPort, p.Number})
	fmt.Printf("adding own port: %v\n", ip+":"+listenOnPort)
	p.peerlist.lock.Unlock()
	go p.receiveFromChannels()
	//Test on localhost
	conn, err := net.Dial("tcp", ip+":"+connectToPort)
	if err != nil {
		fmt.Println("Could not connect peer")
	} else if conn != nil {
		//defer conn.Close()
		//Make the decoder such that we can decode the messages
		dec := gob.NewDecoder(conn)
		enc := gob.NewEncoder(conn)
		p.connections = append(p.connections, enc)
		p.receivePeers(dec)
		p.connectToPeers(ip + ":" + connectToPort)
		go p.handleConnection(dec)
	}
	go p.listenForConnections(totalPeers, ip, listenOnPort)
}
