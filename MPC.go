package main

import (
	config "MPC/Config"
	exp "MPC/Experiment"
	field "MPC/Fields"
	graph "MPC/Graph"
	p "MPC/Party"
	prot "MPC/Protocol"
	"fmt"
)

func main() {
	runExperiments()

}

func runConfig() {
	path := "Config\\configFiles\\" + "kaare" + ".json"
	config := &config.ReadConfig(path)[0]
	field := field.MakeModPrime(1049)
	peer := p.MkPeer(config)
	protocol := prot.MkProtocol(config, field, peer)
	res := protocol.Run()
	fmt.Printf("Got result: %v\n", res)
}

func runExperiments() {
	e := graph.MkExcel("Increment-Mult", "Delay (ms)")
	exp.IncrementMult(e, 0)
}
