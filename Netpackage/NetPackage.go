package netpackage

type NetPackage struct {
	IpPorts []PeerTuple
	Message Message
	Peer    PeerTuple
}

type PeerTuple struct {
	IpPort string
	Number int
}

type Message struct {
	Share     Share
	Signature string
}

type Share struct {
	Value  int64
	Number int64
}
