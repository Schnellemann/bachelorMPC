package party

import (
	aux "MPC/Auxiliary"
	config "MPC/Config"
	netpack "MPC/Netpackage"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"sync"
)

type Peer struct {
	Number       int
	cShare       chan netpack.Share
	connections  connection
	peerlist     *peerList
	wg           *sync.WaitGroup
	config       *config.Config
	decoderMap   map[*gob.Decoder]*ConnectionTuple
	hasSentReady bool
}

type peerList struct {
	hasReceived *sync.WaitGroup
	ipPorts     []string
	lock        sync.Mutex
}

type connection struct {
	c    []*ConnectionTuple
	lock sync.Mutex
}

type ConnectionTuple struct {
	Connection *gob.Encoder
	Number     int
}

func mkPeerList() *peerList {
	pl := new(peerList)
	pl.hasReceived = new(sync.WaitGroup)
	pl.hasReceived.Add(1)
	return pl
}

func MkPeer(config *config.Config) *Peer {
	p := new(Peer)
	p.config = config
	p.Number = int(config.VariableConfig.PartyNr)
	p.peerlist = mkPeerList()
	p.decoderMap = make(map[*gob.Decoder]*ConnectionTuple)
	return p
}

func (p *Peer) StartPeer(shareChannel chan netpack.Share, wg *sync.WaitGroup) {
	p.cShare = shareChannel
	p.wg = wg
	go p.listenForConnections(p.config.VariableConfig.ListenIpPort)
	p.peerlist.lock.Lock()
	p.peerlist.ipPorts = append(p.peerlist.ipPorts, p.config.VariableConfig.ListenIpPort)
	p.peerlist.lock.Unlock()

	if p.config.VariableConfig.ConnectIpPort == "" {
		p.peerlist.hasReceived.Done()
	} else {
		conn, err := net.Dial("tcp", p.config.VariableConfig.ConnectIpPort)
		if err != nil {
			fmt.Printf("Port %v could not connect peer on port %v \n", p.config.VariableConfig.ListenIpPort, p.config.VariableConfig.ConnectIpPort)
		} else if conn != nil {
			dec := gob.NewDecoder(conn)
			enc := gob.NewEncoder(conn)
			//add first connection to map
			conTuble := new(ConnectionTuple)
			conTuble.Connection = enc
			conTuble.Number = p.getPartyNrFromIp(p.config.VariableConfig.ConnectIpPort)
			p.addEntryDecoderMap(dec, conTuble)
			p.connections.lock.Lock()
			p.connections.c = append(p.connections.c, conTuble)
			p.connections.lock.Unlock()
			p.checkReady()
			p.receivePeers(dec)
			p.connectToPeers()
			go p.handleConnection(dec)
			p.broadcastPeer(p.config.VariableConfig.ListenIpPort)
		}
	}

}

func (p *Peer) SendShares(shareList []netpack.Share) {
	p.connections.lock.Lock()
	connections := len(p.connections.c)
	p.connections.lock.Unlock()
	if len(shareList) != connections+1 {
		fmt.Printf("Received sharelist is not of correct length have %v connections, got %v shares, expected %v\n", connections, len(shareList), connections+1)
	}

	p.cShare <- shareList[p.Number-1]
	p.connections.lock.Lock()
	for _, s := range p.connections.c {
		netPackage := new(netpack.NetPackage)
		netPackage.Share = shareList[s.Number-1]
		p.write(s.Connection, netPackage)
	}
	p.connections.lock.Unlock()
}

func (p *Peer) SendShare(share netpack.Share, receiver int) {
	if receiver == p.Number {
		p.cShare <- share
	} else {
		p.connections.lock.Lock()
		for i := range p.connections.c {
			if p.connections.c[i].Number == receiver {
				netPackage := new(netpack.NetPackage)
				netPackage.Share = share
				p.write(p.connections.c[i].Connection, netPackage)
				break
			}
		}
		p.connections.lock.Unlock()
	}

}

func (p *Peer) sendPeerlist(encoder *gob.Encoder) {
	p.peerlist.hasReceived.Wait()
	peerPackage := new(netpack.NetPackage)
	peerPackage.IpPorts = p.peerlist.ipPorts
	p.write(encoder, peerPackage)
}

func (p *Peer) broadcastPeer(ipPort string) {
	newPeerPackage := new(netpack.NetPackage)
	newPeerPackage.Peer = ipPort
	p.connections.lock.Lock()
	for _, c := range p.connections.c {
		p.write(c.Connection, newPeerPackage)
	}
	p.connections.lock.Unlock()
}

func (p *Peer) write(encoder *gob.Encoder, pack *netpack.NetPackage) {
	err := encoder.Encode(pack)
	if err != nil {
		fmt.Println("Could not encode transaction")
		fmt.Println(err)
	}
}

func (p *Peer) SendFinal(share netpack.Share) {
	pack := new(netpack.NetPackage)
	pack.Share = share
	p.connections.lock.Lock()
	for _, conTup := range p.connections.c {
		p.write(conTup.Connection, pack)
	}
	p.connections.lock.Unlock()
}

func (p *Peer) receivePeers(dec *gob.Decoder) {
	recievedPeersPackage := netpack.NetPackage{}
	err := dec.Decode(&recievedPeersPackage)
	if err != nil {
		fmt.Println("Could not decode peer list package from peer")
		return
	}
	p.peerlist.lock.Lock()
	p.peerlist.ipPorts = append(p.peerlist.ipPorts, recievedPeersPackage.IpPorts...)
	p.peerlist.lock.Unlock()
	p.peerlist.hasReceived.Done()
}

func (p *Peer) addConnection(newConnection net.Conn, ip string) {
	encoder := gob.NewEncoder(newConnection)
	decoder := gob.NewDecoder(newConnection)
	conTuble := new(ConnectionTuple)
	if ip != "" {
		conTuble.Number = p.getPartyNrFromIp(ip)
	}
	conTuble.Connection = encoder
	p.addEntryDecoderMap(decoder, conTuble)
	p.connections.lock.Lock()
	p.connections.c = append(p.connections.c, conTuble)
	p.connections.lock.Unlock()
	p.checkReady()
	go p.sendPeerlist(encoder)
	go p.handleConnection(decoder)
}

func (p *Peer) handleConnection(dec *gob.Decoder) {
	//A new peer has connected to us
	//Start receiving packages
	for {
		receivedPackage := &netpack.NetPackage{}
		err := dec.Decode(receivedPackage)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				fmt.Println("Could not decode package from peer")
				fmt.Println(err)
				return
			}
		} else {
			if receivedPackage.Peer != "" {
				p.peerlist.lock.Lock()
				p.peerlist.ipPorts = append(p.peerlist.ipPorts, receivedPackage.Peer)
				p.peerlist.lock.Unlock()
				p.decoderMap[dec].Number = p.getPartyNrFromIp(receivedPackage.Peer)
				p.checkReady()

			} else if receivedPackage.IpPorts == nil {
				s := receivedPackage.Share
				p.cShare <- s
			}
		}
	}
}

func (p *Peer) addEntryDecoderMap(decoder *gob.Decoder, conTuble *ConnectionTuple) {
	p.decoderMap[decoder] = conTuble
}

func (p *Peer) connectToPeers() {
	for _, ip := range p.peerlist.ipPorts {
		//Make sure you don't connect to the initial peer again
		if ip != p.config.VariableConfig.ConnectIpPort && ip != p.config.VariableConfig.ListenIpPort {
			conn, err := net.Dial("tcp", ip)
			if err != nil {
				fmt.Println("Failed to connect to " + ip)
				fmt.Println(err)
			} else {
				p.addConnection(conn, ip)
			}
		}
	}
	p.checkReady()
}

func (p *Peer) listenForConnections(listenOnAddress string) {
	li, err := net.Listen("tcp", listenOnAddress)
	if err != nil {
		fmt.Printf("Error in listening, err: %v\n", err.Error())
		return
	}
	defer li.Close()
	if p.config.VariableConfig.PartyNr == p.config.ConstantConfig.NumberOfParties {
		return
	}
	for {
		conn, err := li.Accept()
		if err != nil {
			fmt.Println("Failed connection on accept")
			return
		}
		p.addConnection(conn, "")
		p.connections.lock.Lock()
		numberOfConnections := len(p.connections.c)
		p.connections.lock.Unlock()
		if numberOfConnections == int(p.config.ConstantConfig.NumberOfParties)-1 {
			break
		}

	}
}

func (p *Peer) getPartyNrFromIp(ip string) int {
	limit := int(p.config.ConstantConfig.NumberOfParties)
	pred := func(i int) bool { return p.config.ConstantConfig.Ipports[i] == ip }
	index := aux.SliceIndex(limit, pred)
	if index == -1 {
		println("Error in finding party number from ip")
	}
	return index + 1
}

func (p *Peer) checkReady() {

	p.connections.lock.Lock()
	defer p.connections.lock.Unlock()
	if p.hasSentReady {
		return
	}
	ready := false
	if len(p.connections.c) == int(p.config.ConstantConfig.NumberOfParties)-1 {
		ready = true
		for _, c := range p.connections.c {
			if c.Number == 0 {
				ready = false
			}
		}
	}
	if ready {
		//End of phase 1
		p.hasSentReady = true
		p.wg.Done()
	}

}
