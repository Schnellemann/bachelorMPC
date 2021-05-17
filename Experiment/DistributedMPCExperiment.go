package experiment

import (
	config "MPC/Config"
	field "MPC/Fields"
	graph "MPC/Graph"
	prot "MPC/Protocol"
	"fmt"
	"strconv"
	"time"
)

func MakeDistributedExperimentFiles(peersPrComputer int, nrOfComputers int, ips []string) {
	for i := 1; i < 11; i++ {
		nrOfParties := nrOfComputers * (peersPrComputer)
		var paths = makePathStrings(nrOfComputers, nrOfParties)
		var secrets = makeRandomSecretList(nrOfParties, 1049)
		exp := makeRandomBalancedMultExpression(nrOfParties, 100)
		confs := config.MakeDistributedConfigs(ips, peersPrComputer, nrOfParties, secrets, exp)
		config.WriteConfig(paths, confs, peersPrComputer)
		peersPrComputer += 3
	}
}

func makePathStrings(numberOfComputers int, numberOfParties int) (paths []string) {
	for i := 0; i < numberOfComputers; i++ {
		paths = append(paths, "com_"+strconv.Itoa(i+1)+"-"+strconv.Itoa(numberOfParties)+"-"+"peers.json")
	}
	return
}

func RunDistributedExperiment(path string, plotter graph.Interface, numberOfParties int) {
	plotter.NewSeries("Number-of-Parties: " + strconv.Itoa(numberOfParties))
	fieldRange := 1049
	configs := config.ReadConfig(path)
	peerlist := getXPeers(configs)
	var channels []chan int64
	var timers []*prot.Times
	for j := 0; j < len(configs); j++ {
		timers = append(timers, new(prot.Times))
	}
	var tProtList []*prot.Times
	for j, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), peerlist[j])
		tprot := prot.MkTimeMeasuringProt(p, c, timers[j])
		tProtList = append(tProtList, tprot.Timer)
		go goProt(tprot, channel)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Done setting up")
	var resultList []int
	for _, c := range channels {
		result := <-c
		fmt.Println(result)
		resultList = append(resultList, int(result))
	}

	if !allSameResults(resultList) {
		fmt.Println("Peers do not agree on the result")
		fmt.Printf("Result: %v \n", resultList)
	}
	avgTProt := prot.AverageTimes(tProtList)
	plotter.AddData(len(configs), avgTProt)
	plotter.Plot()
}
