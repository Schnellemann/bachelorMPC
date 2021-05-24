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

func MakeDistributedIncMults(start int, end int, increment int, ips []string) {
	for i := start; i <= end; i += increment {
		MakeDistributedExperimentFile(3, i, ips, "mults")
	}
}

func MakeDistributedExperimentFile(peersPrComputer int, nrOfMults int, ips []string, iden string) {
	nrOfParties := len(ips) * (peersPrComputer)
	var paths = makePathStrings(len(ips), nrOfMults, iden)
	var secrets = makeRandomSecretList(nrOfParties, 1049)
	exp := makeRandomMultExpression(nrOfParties, nrOfMults)
	confs := config.MakeDistributedConfigs(ips, peersPrComputer, nrOfParties, secrets, exp)
	config.WriteConfig(paths, confs, peersPrComputer)
}

func makePathStrings(numberOfComputers int, variable int, variableIden string) (paths []string) {
	for i := 0; i < numberOfComputers; i++ {
		paths = append(paths, "com_"+strconv.Itoa(i+1)+"-"+strconv.Itoa(variable)+"-"+variableIden+".json")
	}
	return
}

func RunDistributedExperiment(path string, plotter graph.Interface, numberOfMults int) {
	plotter.NewSeries("Number-of-Mults: " + strconv.Itoa(numberOfMults))
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
