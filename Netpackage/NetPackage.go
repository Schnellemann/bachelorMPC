package netpackage

type NetPackage struct {
	IpPorts []string
	Share Share
	Peer    string
}

type PeerTuple struct {
	IpPort string
	Number int
}

type Share struct {
	Value      int64
	Identifier string
}
