package protocol

type Prot interface {
	startNetwork()
	calculate() int64
	setupTree()
	runPreprocess()
	Run() int64
	Destroy()
}
