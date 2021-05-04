package experiment

import (
	config "MPC/Config"
	field "MPC/Fields"
	graph "MPC/Graph"
	party "MPC/Party"
	prot "MPC/Protocol"

	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var ip string = "127.0.1.1"

func getXPeers(configList []*config.Config) []*party.Peer {
	var peers []*party.Peer
	for _, c := range configList {
		peer := party.MkPeer(c)
		peers = append(peers, peer)
	}
	return peers
}

func goProt(prot prot.Prot, result chan int64) {
	res := prot.Run()
	result <- res
}

func allSameResults(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i] != a[0] {
			return false
		}
	}
	return true
}

func makeRandomMultExpression(nrOfPeers int, nrOfMultiplication int) string {
	expression := "p" + strconv.Itoa(rand.Intn(nrOfPeers)+1)
	for i := 0; i < nrOfMultiplication; i++ {
		peerNr := rand.Intn(nrOfPeers) + 1
		expression += "*p" + strconv.Itoa(peerNr)
	}

	return expression
}

func makeRandomSecretList(nrOfParties int, field int) []int {
	var secretList []int
	for i := 0; i < nrOfParties; i++ {
		secret := rand.Intn(field)
		secretList = append(secretList, secret)
	}
	return secretList
}

//=========================================================| Fast Experiments |==============================================================================================

//Increment peers
func IncPeers(plotter graph.Interface) {
	fieldRange := 1049
	for i := 3; i < 70; i += 10 {
		fmt.Printf("Starting Experiment with %v peers. \n", i)
		secretList := makeRandomSecretList(i, fieldRange)
		expression := makeRandomMultExpression(len(secretList), 20)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
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
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), peerlist[j])
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

//Increment multiplication
func IncMult(plotter graph.Interface) {
	fieldRange := 1049
	for i := 1000; i < 20000; i += 2000 {
		fmt.Printf("Starting Experiment with %v multiplication. \n", i)
		secretList := makeRandomSecretList(10, fieldRange)
		expression := makeRandomMultExpression(len(secretList), i)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
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
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), peerlist[j])
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

//Increment bandwidth
func IncBandwidth(plotter graph.Interface) {
	fieldRange := 1049
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
	plotter.Plot()
}

//Increment delay
func IncDelay(plotter graph.Interface) {
	fieldRange := 1049
	for i := 10; i < 100; i += 10 {
		fmt.Printf("Starting Experiment with %v delay \n", i)
		secretList := makeRandomSecretList(10, fieldRange)
		expression := makeRandomMultExpression(len(secretList), i)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		var delayPeerlist []party.IPeer
		//Convert to delayPeer
		for j, p := range peerlist {
			dPeer := party.MkDelayedPeer(configs[j], time.Duration(i), p)
			delayPeerlist = append(delayPeerlist, dPeer)
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
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), delayPeerlist[j])
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
