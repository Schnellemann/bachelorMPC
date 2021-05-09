package protocol

type Prot interface {
	startNetwork()
	calculate() int64
	setupTree()
	Run() int64
}
