package experiment

import (
	config "MPC/Config"
	field "MPC/Fields"
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

func makeRandomExpression(nrOfPeers int, nrOfMultiplication int) string {
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

//Increment peers
func incPeers() {
	fieldRange := 13
	for i := 3; i < 100; i += 10 {
		secretList := makeRandomSecretList(i, fieldRange)
		expression := makeRandomExpression(len(secretList), 20)

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

//Increment add or scalar instructions

//Increment multiplication

//Increment bandwidth

//Increment delay

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
