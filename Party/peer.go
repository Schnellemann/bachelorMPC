package party

import (
	aux "MPC/Auxiliary"
	netpack "MPC/Netpackage"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
)

type Peer struct {
	Number       int
	ipListen     string
	cConnections chan net.Conn
	cPackages    chan *netpack.NetPackage
	cMessages    chan netpack.Message
	connections  []ConnectionTuple
	peerlist     *peerList
}

type peerList struct {
	ipPorts []netpack.PeerTuple
	lock    sync.Mutex
}

type ConnectionTuple struct {
	Connection *gob.Encoder
	Number     int
}

func mkPeerList() *peerList {
	pl := new(peerList)
	return pl
}

func MkPeer(number int, messageChannel chan netpack.Message) *Peer {
	p := new(Peer)
	p.Number = number
	p.cMessages = messageChannel
	p.peerlist = mkPeerList()
	p.cPackages = make(chan *netpack.NetPackage)
	p.cConnections = make(chan net.Conn)

	return p
}

func (p *Peer) handleConnection(dec *gob.Decoder) {
	//A new peer has connected to us
	//Start receiving packages
	for {
		receivedPackage := &netpack.NetPackage{}
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

func (p *Peer) sendShares() {
	//shareList := p.makeShares(int64(len(p.connections)))
	//for i := 0; i < len(shareList); i++ {
	//netPackage := new(NetPackage)
	//netPackage.Message.Share = shareList[i]
	//This should work because connections are sorted in recieveFromChannels
	//p.write(p.connections[i].Connection, netPackage)
	//}
}

func (p *Peer) receiveFromChannels() {
	for {
		select {
		case newConnection := <-p.cConnections:
			encoder := gob.NewEncoder(newConnection)
			decoder := gob.NewDecoder(newConnection)
			fmt.Printf("RemoteAddr for the new connection is: %v \n", newConnection.RemoteAddr)
			fmt.Printf("All availeble connections for party %v: %v", p.Number, p.connections)
			p.connections = append(p.connections, ConnectionTuple{encoder, p.Number})
			//Sort connections by Number
			p.sortConnections()
			//send peers to the new connections
			p.sendPeerlist(encoder)
			go p.handleConnection(decoder)
		case newPackage := <-p.cPackages:
			p.processPackage(newPackage)
		}
	}
}

func (p *Peer) sortConnections() {
	sort.SliceStable(p.connections, func(i, j int) bool {
		return p.connections[i].Number < p.connections[j].Number
	})
}

func (p *Peer) sendPeerlist(encoder *gob.Encoder) {
	peerPackage := new(netpack.NetPackage)
	peerPackage.IpPorts = p.peerlist.ipPorts
	p.write(encoder, peerPackage)
}

func (p *Peer) write(encoder *gob.Encoder, pack *netpack.NetPackage) {
	err := encoder.Encode(pack)
	if err != nil {
		fmt.Println("Could not encode transaction trying again")
		fmt.Println(err)
		p.write(encoder, pack)
	}
}

func (p *Peer) processPackage(pack *netpack.NetPackage) {
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

func (p *Peer) processMessage(message netpack.Message) {

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
	recievedPeersPackage := netpack.NetPackage{}
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
	newPeerPackage := new(netpack.NetPackage)
	newPeerPackage.Peer = netpack.PeerTuple{IpPort: ipPort, Number: p.Number}
	for _, c := range p.connections {
		p.write(c.Connection, newPeerPackage)
	}
}

func (p *Peer) startPeer(totalPeers int, ip string, connectToPort string, listenOnPort string) {
	p.ipListen = ip + ":" + listenOnPort
	p.peerlist.lock.Lock()
	p.peerlist.ipPorts = append(p.peerlist.ipPorts, netpack.PeerTuple{ip + ":" + listenOnPort, p.Number})
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
		p.connections = append(p.connections, ConnectionTuple{enc, p.Number})
		p.receivePeers(dec)
		p.connectToPeers(ip + ":" + connectToPort)
		go p.handleConnection(dec)
	}
	go p.listenForConnections(totalPeers, ip, listenOnPort)
}
