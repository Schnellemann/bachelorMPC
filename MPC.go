package main

import (
	config "MPC/Config"
	exp "MPC/Experiment"
	field "MPC/Fields"
	graph "MPC/Graph"
	p "MPC/Party"
	prot "MPC/Protocol"
	"fmt"
	"strconv"
)

func main() {
	runExperiments()
}

func runConfig() {
	path := "Config\\configFiles\\" + "kaare" + ".json"
	config := config.ReadConfig(path)[0]
	field := field.MakeModPrime(1049)
	peer := p.MkPeer(config)
	protocol := prot.MkProtocol(config, field, peer)
	res := protocol.Run()
	fmt.Printf("Got result: %v\n", res)
}

func runExperiments() {
	var e graph.Interface
	//e = graph.MkExcel("Increment Peer", "Peers")
	//exp.IncrementPeer(e)
	//e = graph.MkExcel("Increment-Mult+Delay", "Number of mults")
	//exp.IncrementMultAndDelay(e)
	//e = graph.MkExcel("Increment-Mult", "Delay (ms)")
	//exp.IncrementMult(e)
	//e = graph.MkPlotter("Increment Delay", "", "png", "Number of Peers")
	//exp.IncDelay(e)
	numberOfMults := 14000
	e = graph.MkExcel("Distributed-Mult-"+strconv.Itoa(numberOfMults), "Mults")
	computerNr := 1

	path := "com_" + strconv.Itoa(computerNr) + "-" + strconv.Itoa(numberOfMults) + "-mults.json"
	exp.RunDistributedExperiment(path, e, numberOfMults)

}

func makeDistributed() {
	var ips = []string{
		"192.168.1.193",
		"192.168.1.141",
		"192.168.1.248",
	}
	exp.MakeDistributedIncMults(1000, 20000, 1000, ips)
}
