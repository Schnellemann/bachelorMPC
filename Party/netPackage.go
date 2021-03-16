package party

type NetPackage struct {
	IpPorts []string
	Message Message
	Peer    string
}

type Message struct {
	TreeLocation TreeLocation
	Share        Share
	Signature    string
}

type TreeLocation struct {
}
