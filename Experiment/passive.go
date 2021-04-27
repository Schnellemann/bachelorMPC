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
	expression := ""
	for i := 0; i < nrOfMultiplication; i++ {
		peerNr := rand.Intn(nrOfPeers) + 1
		expression += "p" + strconv.Itoa(peerNr)
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
func incPeers() {
	fieldRange := 13
	var xyList []graph.XY
	for i := 3; i < 100; i += 10 {
		secretList := makeRandomSecretList(i, fieldRange)
		expression := makeRandomMultExpression(len(secretList), 20)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		var channels []chan int64
		for i, c := range configs {
			channel := make(chan int64)
			channels = append(channels, channel)
			//Make protocol
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), peerlist[i])
			tprot := prot.MkTimeMeasuringProt(p, c)
			go goProt(tprot, channel)
			timeStruct := tprot.Timer
			//Only count calculate and preprocessing for the experiment TODO: maybe some others?
			y := timeStruct.Calculate + timeStruct.Preprocess
			xyList = append(xyList, graph.XY{float64(i), float64(y)})
			//TODO how does Jens want the times????
			time.Sleep(200 * time.Millisecond)
		}

		//Change this so it checks that all the results are similar
		var resultList []int
		for _, c := range channels {
			result := <-c
			resultList = append(resultList, int(result))
		}
		if !allSameResults(resultList) {
			fmt.Println("Peers do not agree on the result")
		}
	}
	graph.PlotGraph("increment Peers", xyList, "Jens er dum", "png")
}

//Increment add or scalar instructions

//Increment multiplication
func incMult() {
	fieldRange := 13
	for i := 20; i < 100; i += 10 {
		secretList := makeRandomSecretList(10, fieldRange)
		expression := makeRandomMultExpression(len(secretList), i)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		var channels []chan int64
		for i, c := range configs {
			channel := make(chan int64)
			channels = append(channels, channel)
			//Make protocol
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), peerlist[i])
			tprot := prot.MkTimeMeasuringProt(p, c)
			go goProt(tprot, channel)
			//timeStruct := tprot.Timer
			//TODO how does Jens want the times????
			time.Sleep(200 * time.Millisecond)
		}

		//Change this so it checks that all the results are similar
		var resultList []int
		for _, c := range channels {
			result := <-c
			resultList = append(resultList, int(result))
		}
		if !allSameResults(resultList) {
			fmt.Println("Peers do not agree on the result")
		}
	}
}

//Increment bandwidth
func incBandwidth() {
	fieldRange := 13
	for i := 10; i < 100; i += 10 {
		secretList := makeRandomSecretList(10, fieldRange)
		expression := makeRandomMultExpression(len(secretList), 20)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		var bandwidthPeerlist []party.IPeer
		//Convert to bandwidthPeer
		for j, p := range peerlist {
			bPeer := party.MkBandwidthPeer(configs[j], p, i, 10*time.Millisecond)
			bandwidthPeerlist = append(bandwidthPeerlist, bPeer)
		}

		var channels []chan int64
		for i, c := range configs {
			channel := make(chan int64)
			channels = append(channels, channel)
			//Make protocol
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), bandwidthPeerlist[i])
			tprot := prot.MkTimeMeasuringProt(p, c)
			go goProt(tprot, channel)
			//timeStruct := tprot.Timer
			//TODO how does Jens want the times????
			time.Sleep(200 * time.Millisecond)
		}

		//Change this so it checks that all the results are similar
		var resultList []int
		for _, c := range channels {
			result := <-c
			resultList = append(resultList, int(result))
		}
		if !allSameResults(resultList) {
			fmt.Println("Peers do not agree on the result")
		}
	}
}

func incBandwidthPunish() {
	fieldRange := 13
	for i := 10; i < 100; i += 10 {
		secretList := makeRandomSecretList(10, fieldRange)
		expression := makeRandomMultExpression(len(secretList), 20)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		var bandwidthPeerlist []party.IPeer
		//Convert to bandwidthPeer
		for j, p := range peerlist {
			bPeer := party.MkBandwidthPeer(configs[j], p, 10, time.Duration(i)*time.Millisecond)
			bandwidthPeerlist = append(bandwidthPeerlist, bPeer)
		}

		var channels []chan int64
		for i, c := range configs {
			channel := make(chan int64)
			channels = append(channels, channel)
			//Make protocol
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), bandwidthPeerlist[i])
			tprot := prot.MkTimeMeasuringProt(p, c)
			go goProt(tprot, channel)
			//timeStruct := tprot.Timer
			//TODO how does Jens want the times????
			time.Sleep(200 * time.Millisecond)
		}

		//Change this so it checks that all the results are similar
		var resultList []int
		for _, c := range channels {
			result := <-c
			resultList = append(resultList, int(result))
		}
		if !allSameResults(resultList) {
			fmt.Println("Peers do not agree on the result")
		}
	}
}

//Increment delay
func incDelay() {
	fieldRange := 13
	//TODO maybe change the increment????
	for i := 200; i < 2000; i += 100 {
		secretList := makeRandomSecretList(10, fieldRange)
		expression := makeRandomMultExpression(len(secretList), 20)

		configs := config.MakeConfigs(ip, expression, secretList)
		peerlist := getXPeers(configs)
		var bandwidthPeerlist []party.IPeer
		//Convert to delayPeer
		for j, p := range peerlist {
			bPeer := party.MkDelayedPeer(configs[j], time.Duration(i)*time.Millisecond, p)
			bandwidthPeerlist = append(bandwidthPeerlist, bPeer)
		}

		var channels []chan int64
		for i, c := range configs {
			channel := make(chan int64)
			channels = append(channels, channel)
			//Make protocol
			p := prot.MkProtocol(c, field.MakeModPrime(int64(fieldRange)), bandwidthPeerlist[i])
			tprot := prot.MkTimeMeasuringProt(p, c)
			go goProt(tprot, channel)
			//timeStruct := tprot.Timer
			//TODO how does Jens want the times????
			time.Sleep(200 * time.Millisecond)
		}

		//Change this so it checks that all the results are similar
		var resultList []int
		for _, c := range channels {
			result := <-c
			resultList = append(resultList, int(result))
		}
		if !allSameResults(resultList) {
			fmt.Println("Peers do not agree on the result")
		}
	}
}

//Increment peers and multiplication

//Increment peers and bandwidth

//Increment peers and delay

//Increment peers, multiplication and bandwidth

//Increment peers, multiplication and delay

//Increment peers, multiplication, bandwidth and delay

//Increment multiplication and bandwidth

//Increment multiplication and delay

//Increment multiplication, bandwidth and delay

//Increment bandwidth and delay

//=========================================================| Slow Experiments |==============================================================================================
