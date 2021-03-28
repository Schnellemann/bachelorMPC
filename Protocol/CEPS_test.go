package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	netpack "MPC/Netpackage"
	party "MPC/Party"
	"strconv"
	"testing"
	"time"
)

var ip string = "127.0.1.1"

var peerlist []*mockPeer

type mockPeer struct {
	ShareChannel chan netpack.Share
	finalSend    netpack.Share
	partyNr      int
}

func mkMockPeer(partyNr int) *mockPeer {
	m := new(mockPeer)
	m.partyNr = partyNr
	return m

}

func (m *mockPeer) StartPeer(shareChannel chan netpack.Share) {
	m.ShareChannel = shareChannel
}

func (m *mockPeer) SendShares(shares []netpack.Share) {
	for j := 0; j < len(peerlist); j++ {
		if j != m.partyNr {
			msgToSend := shares[j]
			peerlist[j].ShareChannel <- msgToSend
		}
	}
}

func (m *mockPeer) SetProgress(progress chan int) {
	progress <- 1
}

func (m *mockPeer) SendFinal(share netpack.Share) {
	for j := 0; j < len(peerlist); j++ {
		if j != m.partyNr {
			peerlist[j].ShareChannel <- share
		}
	}
}

func TestOutputReconstruction(t *testing.T) {
	//Start up the mock
	shareChannel := make(chan netpack.Share)
	mockPeer := mkMockPeer(1)
	mockPeer.StartPeer(shareChannel)

	//Config
	vconfig := config.VariableConfig{PartyNr: 1, Secret: 5}
	cconfig := config.ConstantConfig{Expression: "p1", NumberOfParties: 5}
	config := config.Config{VariableConfig: vconfig, ConstantConfig: cconfig}

	theField := field.MakeModPrime(11)
	prot := mkProtocol(&config, theField, mockPeer)
	output := prot.run()

	if output != config.VariableConfig.Secret {
		t.Errorf("Outputreconstruction does not work correctly, expected %v but got %v", config.VariableConfig.Secret, output)
	}

}

func makeConfigs(expression string, secrets []int) []*config.Config {
	var configList []*config.Config
	var listenIpPorts []string
	var connectIpPorts []string
	for i := 0; i < len(secrets); i++ {
		ListenPort := 40000 + i*10
		listenIpPorts = append(listenIpPorts, (ip + ":" + strconv.Itoa(ListenPort)))
		var connectToIpPort string
		if i == 0 {
			connectToIpPort = ""
		} else {
			connectToIpPort = (ip + ":" + strconv.Itoa(40000+(i-1)*10))
		}
		connectIpPorts = append(connectIpPorts, connectToIpPort)
	}

	for i, s := range secrets {

		vconfig := config.VariableConfig{PartyNr: float64(i + 1), Secret: int64(s), ListenIpPort: listenIpPorts[i], ConnectIpPort: connectIpPorts[i]}
		cconfig := config.ConstantConfig{Expression: expression, NumberOfParties: float64(len(secrets)), Ipports: listenIpPorts}
		config := config.Config{VariableConfig: vconfig, ConstantConfig: cconfig}
		configList = append(configList, &config)
	}
	return configList
}

func getXMockPeers(numberOfPeers int) []*mockPeer {
	var mocks []*mockPeer
	for i := 0; i < numberOfPeers; i++ {
		mock := mkMockPeer(i + 1)
		mocks = append(mocks, mock)
	}
	return mocks
}

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
	configs := makeConfigs("p1+p2", []int{4, 7, 3, 2, 1})
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

	configs := makeConfigs("2*p2", []int{4, 7, 3, 2})
	peerlist = getXMockPeers(4)
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
		if result != 1 {
			t.Errorf("Scalar does not work correctly peer %v expected %v but got %v", i+1, 1, result)
		}
	}

}

func TestMultiply(t *testing.T) {

}

func TestCombined(t *testing.T) {

}
