package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	netpack "MPC/Netpackage"
	"testing"
)

type mockPeer struct {
}

func mkMockPeer() *mockPeer {
	m := new(mockPeer)
	return m

}

func (m *mockPeer) StartPeer(shareChannel chan netpack.Share) {
}

func (m *mockPeer) SendShares(shares []netpack.Share) {
}

func (m *mockPeer) SetProgress(progress chan int) {
	progress <- 1
}

func (m *mockPeer) SendFinal(share netpack.Share) {
}

func (m *mockPeer) setDone() {

}

func TestOutputReconstruction(t *testing.T) {
	//Start up the mock
	shareChannel := make(chan netpack.Share)
	mockPeer := mkMockPeer()
	mockPeer.StartPeer(shareChannel)

	//Config
	vconfig := config.VariableConfig{PartyNr: 1, Secret: 5}
	cconfig := config.ConstantConfig{Expression: "p1", NumberOfParties: 5}
	config := config.Config{VariableConfig: vconfig, ConstantConfig: cconfig}

	theField := field.MakeModPrime(11)
	prot := mkProtocol(&config, config.VariableConfig.Secret, theField, mockPeer)

	output := prot.run()

	if output != config.VariableConfig.Secret {
		t.Errorf("Outputreconstruction does not work correctly, expected %v but got %v", config.VariableConfig.Secret, output)
	}

}

func TestAdd(t *testing.T) {

}

func TestScalar(t *testing.T) {

}

func TestMultiply(t *testing.T) {

}

func TestCombined(t *testing.T) {

}
