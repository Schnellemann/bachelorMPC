package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	party "MPC/Party"
	"testing"
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

func (prot *Ceps) goProt(result chan int64) {
	res := prot.run()
	result <- res
}

func TestAdd(t *testing.T) {
	configs := config.MakeConfigs(ip, "p1+p2", []int{4, 7, 3})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(13), peerlist[i])
		go prot.goProt(channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 11 {
			t.Errorf("Addition does not work correctly peer %v expected %v but got %v", i+1, 11, result)
		}
	}

}

func TestScalar(t *testing.T) {

	configs := config.MakeConfigs(ip, "2*p3", []int{4, 7, 3, 2, 1})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(13), peerlist[i])
		go prot.goProt(channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 6 {
			t.Errorf("Scalar does not work correctly peer %v expected %v but got %v", i+1, 6, result)
		}
	}

}

func TestMultiply(t *testing.T) {
	configs := config.MakeConfigs(ip, "p2*p3", []int{4, 7, 3})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(13), peerlist[i])
		go prot.goProt(channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 8 {
			t.Errorf("Multiply does not work correctly peer %v expected %v but got %v", i+1, 8, result)
		}
	}
}

func TestCombined(t *testing.T) {

}
