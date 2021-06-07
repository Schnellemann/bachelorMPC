package main

import (
	config "MPC/Config"
	exp "MPC/Experiment"
	field "MPC/Fields"
	graph "MPC/Graph"
	p "MPC/Party"
	prot "MPC/Protocol"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if args[0] == "" {
		fmt.Println("Running Experiment")
		runExperiments()
	}
	if args[0] == "files" {
		makeDistributed()
	}
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
	/*
		Increasing peer
	*/
	e = graph.MkExcel("Increment Peer", "Peers")
	exp.IncrementPeer(e)
	/*
		Increasing amount of multiplications and peer delay
	*/
	//e = graph.MkExcel("Increment-Mult+Delay", "Number of mults")
	//exp.IncrementMultAndDelay(e)

	/*
		Increment mult and bandwidth size
	*/
	//e = graph.MkExcel("Increment-Mult+Bandwidth", "Number of mults")
	//exp.IncrementMultAndBandwidth(e)

	/*
		Increasing amount of multiplications
	*/
	//e = graph.MkExcel("Increment-Mult", "Delay (ms)")
	//exp.IncrementMult(e)

	/*
		Increasing message delay
	*/
	//e = graph.MkPlotter("Increment Delay", "", "png", "Number of Peers")
	//exp.IncDelay(e)

	/*
		Distributed experiment given a path to a config.json file
	*/
	// numberOfMults := 10000
	// e = graph.MkExcel("Distributed-Mult-"+strconv.Itoa(numberOfMults), "Mults")
	// computerNr := 1

	// path := "com_" + strconv.Itoa(computerNr) + "-" + strconv.Itoa(numberOfMults) + "-mults.json"
	// exp.RunDistributedExperiment(path, e, numberOfMults)

}

func makeDistributed() {
	/*
		Insert Ip's here to create config files for distributed experiment
	*/
	var ips = []string{}
	exp.MakeDistributedIncMults(1000, 20000, 1000, ips)
}
