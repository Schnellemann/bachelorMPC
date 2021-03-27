package party

import (
	aux "MPC/Auxiliary"
	config "MPC/Config"
	netpack "MPC/Netpackage"
	"encoding/gob"
	"fmt"
	"net"
	"sort"
	"sync"
)

type Peer struct {
	Number      int
	cShare      chan netpack.Share
	connections []ConnectionTuple
	peerlist    *peerList
	progress    chan int
	config      *config.Config
	decoderMap  map[*gob.Decoder]*ConnectionTuple
}

type peerList struct {
	ipPorts []string
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

func (p *Peer) SetProgress(progress chan int) {
	p.progress = progress
}

func MkPeer(config *config.Config) *Peer {
	p := new(Peer)
	p.config = config
	p.Number = int(config.VariableConfig.PartyNr)
	p.peerlist = mkPeerList()
	p.decoderMap = make(map[*gob.Decoder]*ConnectionTuple)
	return p
}

func (p *Peer) StartPeer(shareChannel chan netpack.Share) {
	p.cShare = shareChannel
	p.peerlist.lock.Lock()
	p.peerlist.ipPorts = append(p.peerlist.ipPorts, p.config.VariableConfig.ListenIpPort)
	p.peerlist.lock.Unlock()
	//Test on localhost
	if p.config.VariableConfig.ConnectIpPort != "" {
		conn, err := net.Dial("tcp", p.config.VariableConfig.ConnectIpPort)
		if err != nil {
			fmt.Println("Could not connect peer")
		} else if conn != nil {
			//defer conn.Close()
			//Make the decoder such that we can decode the messages
			dec := gob.NewDecoder(conn)
			enc := gob.NewEncoder(conn)
			conTuble := new(ConnectionTuple)
			conTuble.Connection = enc
			conTuble.Number = p.getPartyNrFromIp(p.config.VariableConfig.ConnectIpPort)
			p.addEntryDecoderMap(dec, conTuble)
			p.connections = append(p.connections, *conTuble)
			p.receivePeers(dec)
			p.connectToPeers()
			go p.handleConnection(dec)
			p.broadcastPeer(p.config.VariableConfig.ListenIpPort)
		}
	}
	go p.listenForConnections(int(p.config.ConstantConfig.NumberOfParties), p.config.VariableConfig.ListenIpPort)
}

//Send Methods
func (p *Peer) SendShares(shareList []netpack.Share) {
	for i := 0; i < len(shareList); i++ {
		netPackage := new(netpack.NetPackage)
		netPackage.Share = shareList[i]
		//This should work because connections are sorted in recieveFromChannels
		p.write(p.connections[i].Connection, netPackage)
	}
}

func (p *Peer) sendPeerlist(encoder *gob.Encoder) {
	peerPackage := new(netpack.NetPackage)
	peerPackage.IpPorts = p.peerlist.ipPorts
	p.write(encoder, peerPackage)
}

func (p *Peer) broadcastPeer(ipPort string) {
	newPeerPackage := new(netpack.NetPackage)
	newPeerPackage.Peer = ipPort
	for _, c := range p.connections {
		p.write(c.Connection, newPeerPackage)
	}
}

func (p *Peer) write(encoder *gob.Encoder, pack *netpack.NetPackage) {
	err := encoder.Encode(pack)
	if err != nil {
		fmt.Println("Could not encode transaction trying again")
		fmt.Println(err)
		p.write(encoder, pack)
	}
}

func (p *Peer) SendFinal(share netpack.Share) {
	pack := new(netpack.NetPackage)
	pack.Share = share
	for _, conTup := range p.connections {
		p.write(conTup.Connection, pack)
	}
}

// Recieve Methods
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
	p.connections = append(p.connections, *conTuble)
	//send peers to the new connections
	p.sendPeerlist(encoder)
	go p.handleConnection(decoder)
}

// Internal functions
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
			if receivedPackage.Peer != "" {
				//This is a peer broadcast, so add the peer to the peer list and connect the encoder to the p.number.
				p.peerlist.lock.Lock()
				p.peerlist.ipPorts = append(p.peerlist.ipPorts, receivedPackage.Peer)
				p.peerlist.lock.Unlock()
				p.decoderMap[dec].Number = p.getPartyNrFromIp(receivedPackage.Peer)

			} else if receivedPackage.IpPorts == nil {
				//If we receive IpPorts we should ignore it, we handle this
				//Only if we're actually waiting for the peerlist
				s := receivedPackage.Share
				p.cShare <- s
			}
		}
	}
}

func (p *Peer) sortConnections() {
	sort.SliceStable(p.connections, func(i, j int) bool {
		return p.connections[i].Number < p.connections[j].Number
	})
}

func (p *Peer) addEntryDecoderMap(decoder *gob.Decoder, conTuble *ConnectionTuple) {
	p.decoderMap[decoder] = conTuble
}

func (p *Peer) connectToPeers() {
	fmt.Printf("Peerlist at connect to Peers %v\n", p.peerlist.ipPorts)
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
}

func (p *Peer) listenForConnections(totalPeers int, listenOnAddress string) {
	li, err := net.Listen("tcp", listenOnAddress)
	if err != nil {
		fmt.Println("Error in listening")
		return
	}
	defer li.Close()
	fmt.Println("Other peers can connect to me on the following ip:port")
	fmt.Println("Address " + ": " + p.config.VariableConfig.ListenIpPort)
	i := len(p.connections) + 1
	for i < int(p.config.ConstantConfig.NumberOfParties) {
		conn, err := li.Accept()
		//TODO update Connectionstuple with encoder and number of the connected peer
		if err != nil {
			fmt.Println("Failed connection on accept")
			return
		}
		p.addConnection(conn, "")
		i++
	}
	//End of phase 1
	p.progress <- 1
}

func (p *Peer) getPartyNrFromIp(ip string) int {
	return aux.SliceIndex(int(p.config.ConstantConfig.NumberOfParties-1), func(i int) bool { return p.config.ConstantConfig.Ipports[i] == ip })
}
