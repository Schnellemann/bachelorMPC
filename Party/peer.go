package party

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	Number       int
	cConnections chan net.Conn
	cPackages    chan *NetPackage
	connections  []*gob.Encoder
	peerlist     *peerList
}

type peerList struct {
	ipPorts []string
	lock    sync.Mutex
}

func mkPeerList() *peerList {
	pl := new(peerList)
	return pl
}

func mkPeer() *Peer {
	p := new(Peer)
	p.peerlist = mkPeerList()
	p.cPackages = make(chan *NetPackage)
	p.cConnections = make(chan net.Conn)

	return p
}

func startPeer(ip string, connectToPort string, listenOnPort string) {
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
	go p.listenForConnections(ip, listenOnPort)
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
			p.sendPeers(encoder)
			go p.handleConnection(decoder)

		case newPackage := <-p.cPackages:
			p.processPackage(newPackage)
		}
	}
}

func (p *Peer) sendPeers(encoder *gob.Encoder) {
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
	if pack.Peer != "" {
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
	for _, ip := range p.peerlist.ipPorts {
		//Make sure you don't connect to the initial peer again
		if ip != initialConn {
			conn, err := net.Dial("tcp", ip)
			if err != nil {
				fmt.Println("Failed to connect to " + ip)
				fmt.Println(err)
			} else {
				p.cConnections <- conn
			}
		}
	}
}

func (p *Peer) startPeer(ip string, connectToPort string, listenOnPort string) {
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
	go p.listenForConnections(ip, listenOnPort)

}
