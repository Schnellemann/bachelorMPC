package experiment

import (
	config "MPC/Config"
	field "MPC/Fields"
	graph "MPC/Graph"
	prot "MPC/Protocol"

	"fmt"
	"time"
)

//Increment multiplication
//Uses 10 peers
func IncrementMult(plotter graph.Interface) {
	incrementMultWithDelay(plotter, 0, 2000, 50000, 2000)
	plotter.Plot()
}

func incrementMultWithDelay(plotter graph.Interface, delay time.Duration, start int, end int, increment int) {
	fieldRange := 1049
	plotter.NewSeries("Mult with delay " + delay.String())
	for i := start; i <= end; i += increment {
		fmt.Printf("Starting Experiment with %v multiplication. \n", i)
		secretList := makeRandomSecretList(10, fieldRange)
		expression := makeRandomMultExpression(len(secretList), i)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		delayedPeers := getDelayedPeers(configs, peerlist, delay)
		var channels []chan int64
		var timers []*prot.Times
		for j := 0; j < i; j++ {
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