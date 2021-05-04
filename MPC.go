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

func runExperiments() {
	e := graph.MkPlotter("Increment Delay", "", "png", "Number of Peers")
	exp.IncDelay(e)

}

func runConfig() {
	path := "jens" + ".json"
	config := &config.ReadConfig(path)[0]
	field := field.MakeModPrime(1049)
	peer := p.MkPeer(config)
	protocol := prot.MkProtocol(config, field, peer)
	res := protocol.Run()
	fmt.Printf("Got result: %v\n", res)
}
