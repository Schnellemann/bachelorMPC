package experiment

import (
	config "MPC/Config"
	field "MPC/Fields"
	graph "MPC/Graph"
	prot "MPC/Protocol"

	"fmt"
	"time"
)

func IncrementPeer(plotter graph.Interface) {
	incrementPeerWithDelay(plotter, 0, 10, 90, 10, makeRandomMultExpression)
	plotter.Plot()
}

func incrementPeerWithDelay(plotter graph.Interface, delay time.Duration, start int, end int, increment int, multStrat randomMultMaker) {
	fieldRange := 1049
	plotter.NewSeries("Peers with delay " + delay.String())
	for i := start; i <= end; i += increment {
		fmt.Printf("Starting Experiment with %v peers. \n", i)
		secretList := makeRandomSecretList(i, fieldRange)
		expression := multStrat(len(secretList), 20)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		delayedPeers := getDelayedPeers(configs, peerlist, delay)
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
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), delayedPeers[j])
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
}
