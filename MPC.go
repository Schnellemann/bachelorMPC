package main

import (
	exp "MPC/Experiment"
	graph "MPC/Graph"
)

func main() {
	e := graph.MkExcel()
	exp.IncDelay(e)
	e = graph.MkExcel()
	exp.IncBandwidth(e)
	e = graph.MkExcel()
	exp.IncMult(e)
	e = graph.MkExcel()
	exp.IncPeers(e)
}
