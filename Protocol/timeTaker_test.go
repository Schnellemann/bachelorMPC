package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	"testing"
	"time"
)

func TestTimeTaker(t *testing.T) {
	//l := log.New(os.Stdout, "TestTimeTaker> ", log.LstdFlags)
	configs := config.MakeConfigs(ip, "((p1*p1)*((p1*p1)*(p2*p2)))*((p2*p2)*(p3*p3))", []int{2, 3, 5})
	peerlist := getXPeers(configs)
	var channels []chan int64
	for i, c := range configs {
		channel := make(chan int64)
		channels = append(channels, channel)
		//Make protocol
		prot := MkProtocol(c, field.MakeModPrime(43), peerlist[i])
		tprot := MkTimeMeasuringProt(prot, c)
		go goProt(tprot, channel)
		time.Sleep(200 * time.Millisecond)
	}
	for i, c := range channels {
		result := <-c
		if result != 21 {
			t.Errorf("Addition does not work correctly peer %v expected %v but got %v", i+1, 21, result)
		}
	}

}
