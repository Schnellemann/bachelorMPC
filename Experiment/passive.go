package experiment

import (
	config "MPC/Config"
	party "MPC/Party"
	prot "MPC/Protocol"
	"math/rand"
	"strconv"
	"time"
)

var ip string = "127.0.1.1"

func getXPeers(configList []*config.Config) []party.IPeer {
	var peers []party.IPeer
	for _, c := range configList {
		peer := party.MkPeer(c)
		peers = append(peers, peer)
	}
	return peers
}

func getDelayedPeers(configs []*config.Config, peerlist []party.IPeer, delay time.Duration) []party.IPeer {
	var delayPeerlist []party.IPeer
	//Convert to delayPeer
	for j, p := range peerlist {
		dPeer := party.MkDelayedPeer(configs[j], delay, p)
		delayPeerlist = append(delayPeerlist, dPeer)
	}
	return delayPeerlist
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

type randomMultMaker func(int, int) string

func makeRandomMultExpression(nrOfPeers int, nrOfMultiplication int) string {
	expression := "p" + strconv.Itoa(rand.Intn(nrOfPeers)+1)
	for i := 0; i < nrOfMultiplication; i++ {
		peerNr := rand.Intn(nrOfPeers) + 1
		expression += "*p" + strconv.Itoa(peerNr)
	}

	return expression
}

func makeRandomBalancedMultExpression(nrOfPeers int, nrOfMultiplication int) string {
	if nrOfMultiplication > 2 {
		leftNr := nrOfMultiplication / 2
		left := makeRandomBalancedMultExpression(nrOfPeers, leftNr)
		right := makeRandomBalancedMultExpression(nrOfPeers, nrOfMultiplication-leftNr)
		return "(" + left + "*" + right + ")"
	}
	peerNr1 := strconv.Itoa(rand.Intn(nrOfPeers) + 1)
	if nrOfMultiplication == 2 {
		peerNr2 := strconv.Itoa(rand.Intn(nrOfPeers) + 1)
		return "(p" + peerNr1 + "*p" + peerNr2 + ")"
	} else {
		return "(p" + peerNr1 + ")"
	}
}

func makeRandomSecretList(nrOfParties int, field int) []int {
	var secretList []int
	for i := 0; i < nrOfParties; i++ {
		secret := rand.Intn(field)
		secretList = append(secretList, secret)
	}
	return secretList
}
