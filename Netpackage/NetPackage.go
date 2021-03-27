package netpackage

type NetPackage struct {
	IpPorts []string
	Message Message
	Peer    string
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
	Value      int64
	Identifier string
}
