package party

type NetPackage struct {
	IpPorts []PeerTuple
	Message Message
	Peer    PeerTuple
}

type Message struct {
	TreeLocation TreeLocation
	Share        Share
	Signature    string
}

type TreeLocation struct {
}
