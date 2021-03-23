package config

import (
	"testing"
)

var path string = "configTest.json"

func TestConfigLoad(t *testing.T) {
	expString := "(p1 + p2) + p3"
	ipString := "127.0.1.1"
	portString := "6969"
	var conf Config = readConfig(path).Configs[0]
	if conf.Expression != expString {
		t.Errorf("Expression should have been %v, but is %v", expString, conf.Expression)
	}
	if conf.NumberOfParties != 3 {
		t.Errorf("NumberOfParties should have been %v, but is %v", 3, conf.NumberOfParties)
	}
	if conf.Ip != ipString {
		t.Errorf("Ip should have been %v, but is %v", ipString, conf.Ip)
	}
	if conf.Port != portString {
		t.Errorf("Port should have been %v, but is %v", portString, conf.Port)
	}
	if conf.Port != portString {
		t.Errorf("Port should have been %v, but is %v", 1, conf.PartyNr)
	}
}
