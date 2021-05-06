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

func IncBandwidthAndPeers(plotter graph.Interface) {
	fieldRange := 1049
	for i := 3; i < 60; i += 10 {
		fmt.Printf("Starting Series with %v peers. \n", i)
		plotter.NewSeries(strconv.Itoa(i) + " peers")
		for i := 20; i < 200; i *= 2 {
			fmt.Printf("Starting Experiment with %v width. \n", i)
			secretList := makeRandomSecretList(10, fieldRange)
			expression := makeRandomMultExpression(len(secretList), i)

			configs := config.MakeConfigs(ip, expression, secretList)
			peerlist := getXPeers(configs)
			var bandwidthPeerlist []party.IPeer
			//Convert to bandwidthPeer
			for j, p := range peerlist {
				bPeer := party.MkBandwidthPeer(configs[j], p, i, 2*time.Millisecond)
				bandwidthPeerlist = append(bandwidthPeerlist, bPeer)
			}
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
	}
	plotter.Plot()
}
