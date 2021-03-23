package protocol

import (
	config "MPC/Config"
	field "MPC/Fields"
	"testing"
)

func TestMkProtocol(t *testing.T) {
	config := &config.ReadConfig("../Config/ConfigFiles/CEPSConfig.json").Configs[0]
	field := field.MakeModPrime(11)
	proc := mkProtocol(config, 2, field)

	if proc.shamir.degree != 1 {
		t.Errorf("Degree is incorrect there are 1 allowed adversary, but got %v", proc.shamir.degree)
	}
}
