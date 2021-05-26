package experiment

import (
	config "MPC/Config"
	field "MPC/Fields"
	graph "MPC/Graph"
	party "MPC/Party"
	prot "MPC/Protocol"
	"fmt"
	"strconv"
	"time"
)

func IncrementMultAndBandwidth(plotter graph.Interface) {
	for _, bandwidth := range []int{10, 20, 40, 80, 160} {
		incrementMultWithBandwidth(plotter, bandwidth, 1000, 50000, 2000, makeRandomMultExpression, 2*time.Minute)
	}
	plotter.Plot()
}

func incrementMultWithBandwidth(plotter graph.Interface, bandwidth int, start int, end int, increment int, multStrat randomMultMaker, maxExperimentTime time.Duration) {
	fieldRange := 1049
	plotter.NewSeries("Mult with bandwidth " + strconv.Itoa(bandwidth) + "bytes")
	for mults := start; mults <= end; mults += increment {
		fmt.Printf("Starting Experiment with %v multiplication and %v bandwidth. \n", mults, bandwidth)
		secretList := makeRandomSecretList(10, fieldRange)
		expression := multStrat(len(secretList), mults)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		bandwidthPeers := getBandwidthPeers(configs, peerlist, bandwidth)
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
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), bandwidthPeers[j])
			tprot := prot.MkTimeMeasuringProt(p, c, timers[j])
			tProtList = append(tProtList, tprot.Timer)
			go goProt(tprot, channel)
			time.Sleep(100 * time.Millisecond)
		}
		var resultList []int
		for _, c := range channels {
			result := <-c
			resultList = append(resultList, int(result))
		}

		if !allSameResults(resultList) {
			fmt.Println("Peers do not agree on the result")
			fmt.Printf("Result: %v \n", resultList)
		}
		avgTProt := prot.AverageTimes(tProtList)
		plotter.AddData(mults, avgTProt)
		if maxExperimentTime != 0 && avgTProt.Calculate+avgTProt.SetupTree > maxExperimentTime {
			break
		}
	}
}

func getBandwidthPeers(configs []*config.Config, peerlist []party.IPeer, bandwidth int) []party.IPeer {
	var bandwidthPeerlist []party.IPeer
	//Convert to bandwidthPeer
	for j, p := range peerlist {
		bPeer := party.MkBandwidthPeer(configs[j], p, bandwidth, 2*time.Millisecond)
		bandwidthPeerlist = append(bandwidthPeerlist, bPeer)
	}
	return bandwidthPeerlist
}

func IncBandwidth(plotter graph.Interface) {
	plotter.NewSeries("Bandwidth")
	fieldRange := 1049
	for i := 5; i < 400; i *= 2 {
		fmt.Printf("Starting Experiment with %v width. \n", i)
		secretList := makeRandomSecretList(10, fieldRange)
		expression := makeRandomMultExpression(len(secretList), 20)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		bandwidthPeerlist := getBandwidthPeers(configs, peerlist, i)

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
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), bandwidthPeerlist[j])
			tprot := prot.MkTimeMeasuringProt(p, c, timers[j])
			tProtList = append(tProtList, tprot.Timer)
			go goProt(tprot, channel)
			time.Sleep(100 * time.Millisecond)
		}

		//Change this so it checks that all the results are similar
		var resultList []int
		for _, c := range channels {
			result := <-c
			resultList = append(resultList, int(result))
		}

		if !allSameResults(resultList) {
			fmt.Println("Peers do not agree on the result")
			fmt.Printf("Result: %v \n", resultList)
		}
		avgTProt := prot.AverageTimes(tProtList)
		plotter.AddData(i, avgTProt)
	}
	plotter.Plot()
}
