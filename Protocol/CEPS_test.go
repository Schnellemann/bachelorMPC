package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	party "MPC/Party"
	"fmt"
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
		prot := MkProtocol(c, field.MakeModPrime(13), peerlist[i])
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
		prot := MkProtocol(c, field.MakeModPrime(13), peerlist[i])
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
		prot := MkProtocol(c, field.MakeModPrime(13), peerlist[i])
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
		prot := MkProtocol(c, field.MakeModPrime(13), peerlist[i])
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
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
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
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
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

func TestLeftSubTree(t *testing.T) {
	configs := config.MakeConfigs(ip, "(((p1+p1)+p1)+p2)", []int{2, 3, 5})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 9 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 9, result)
		}
	}
}

func TestRightSubTree(t *testing.T) {
	configs := config.MakeConfigs(ip, "(p1+(p1+(p1+p2)))", []int{2, 3, 5})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 9 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 9, result)
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
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
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

func TestLargeBalanced(t *testing.T) {
	configs := config.MakeConfigs(ip,
		"(((((p5)+(p2+p3))+((p10+p9)+(p5+p2)))+(((p6+p8)+(p7+p6))+((p7+p9)+(p9+p8))))+((((p8)+(p8+p9))+((p1+p6)+(p2+p9)))+(((p8+p2)+(p10+p7))+((p8+p2)+(p6+p7)))))",
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := MkProtocol(c, field.MakeModPrime(1049), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Done Setting up the protocols")
	for i, c := range channels {
		result := <-c
		if result != 189 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 189, result)
		}
	}
}

func TestLargeBalancedWith3Peers(t *testing.T) {
	configs := config.MakeConfigs(ip,
		"(((((p1)+(p2+p3))+((p1+p2)+(p1+p2)))+(((p1+p2)+(p3+p3))+((p1+p3)+(p1+p3))))+((((p2)+(p1+p3))+((p1+p2)+(p2+p1)))+(((p3+p2)+(p1+p3))+((p3+p2)+(p2+p3)))))",
		[]int{1, 2, 3})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := MkProtocol(c, field.MakeModPrime(1049), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Done Setting up the protocols")
	for i, c := range channels {
		result := <-c
		if result != 60 {
			t.Errorf("Combined does not work correctly peer %v expected %v but got %v", i+1, 60, result)
		}
	}
}

func TestMultipleMult2(t *testing.T) {
	configs := config.MakeConfigs(ip, "p1*p1*p1*p1*p2*p2*p2*p2*p3*p3", []int{2, 3, 5})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
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
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
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
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
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

func Test33Mult(t *testing.T) {
	//
	var secrets = []int{510, 691, 545, 1005, 911, 194, 675, 803, 195, 162, 233, 629, 204, 990, 159, 709, 845, 285,
		736, 344, 300, 394, 1030, 1012, 799, 439, 921, 936, 158, 500, 990, 384, 278}
	// 158 * 285 * 921 * 510 * 162 * 799 * 159 * 439 * 691 * 285 * 500 * 921 * 921 * 285 * 500 * 158 * 1030 * 911 * 384 * 1012 * 691
	configs := config.MakeConfigs(ip, "p29*p18*p27*p1*p10*p25*p15*p26*p2*p18*p30*p27*p27*p18*p30*p29*p23*p5*p32*p24*p2", secrets)
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := MkProtocol(c, field.MakeModPrime(1049), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 298 {
			t.Errorf("Wrong result expected %v, but got %v at party: %v \n", 298, result, i+1)
		}
	}
}

func Test33PeersAdd(t *testing.T) {
	//
	var secrets = []int{510, 691, 545, 1005, 911, 194, 675, 803, 195, 162, 233, 629, 204, 990, 159, 709, 845, 285,
		736, 344, 300, 394, 1030, 1012, 799, 439, 921, 936, 158, 500, 990, 384, 278}
	// 158 * 285 * 921 * 510 * 162 * 799 * 159 * 439 * 691 * 285 * 500 * 921 * 921 * 285 * 500 * 158 * 1030 * 911 * 384 * 1012 * 691
	configs := config.MakeConfigs(ip, "p29+p18+p27+p1+p10+p25+p15+p26+p2+p18+p30+p27+p27+p18+p30+p29+p23+p5+p32+p24+p2", secrets)
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := MkProtocol(c, field.MakeModPrime(1049), peerlist[i])
		go goProt(prot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 183 {
			t.Errorf("Wrong result expected %v, but got %v at party: %v \n", 183, result, i+1)
		}
	}
}

// p29*p18*p27*p1*p10*p25*p15*p26*p2*p18*p30*p27*p27*p18*p30*p29*p23*p5*p32*p24*p2
// 510 691 545 1005 911 194 675 803 195 162 233 629 204 990 159 709 845 285 736 344 300 394 1030 1012 799 439 921 936 158 500 990 384 278
