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

func goProt(prot Prot, result chan int64) {
	res := prot.Run()
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
		go goProt(prot, channel)
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
		go goProt(prot, channel)
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
	configs := config.MakeConfigs(ip, "p2*p3", []int{11, 4, 9, 2, 5, 7, 5, 3, 6, 11})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(13), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 10 {
			t.Errorf("Multiply does not work correctly peer %v expected %v but got %v", i+1, 10, result)
		}
	}
}

func TestSmallMultiply(t *testing.T) {
	configs := config.MakeConfigs(ip, "p2*p3", []int{11, 4, 9})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(13), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 10 {
			t.Errorf("Multiply does not work correctly peer %v expected %v but got %v", i+1, 10, result)
		}
	}
}

func TestCombined(t *testing.T) {
	//19+12*4+2*40 mod 43 = 18
	configs := config.MakeConfigs(ip, "p1+p2*p4+2*p3", []int{19, 12, 40, 4})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(43), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 18 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 18, result)
		}
	}
}

func TestMultipleAdd(t *testing.T) {
	//19+12+4+40 mod 43 = 32
	configs := config.MakeConfigs(ip, "p1+p2+p4+p3", []int{19, 12, 40, 4})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(43), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 32 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 32, result)
		}
	}
}

func TestMultipleMult(t *testing.T) {
	configs := config.MakeConfigs(ip, "((p1*p1)*((p1*p1)*(p2*p2)))*((p2*p2)*(p3*p3))", []int{2, 3, 5})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(43), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 21 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 21, result)
		}
	}
}

func TestMultSame(t *testing.T) {
	configs := config.MakeConfigs(ip, "p1*p1", []int{2, 3, 5})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(43), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 4 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 4, result)
		}
	}
}

func TestMultipleMultSame(t *testing.T) {
	configs := config.MakeConfigs(ip, "(p1*p1)*(p1*p1)", []int{2, 3, 5})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := mkProtocol(c, field.MakeModPrime(43), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 16 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 16, result)
		}
	}
}
