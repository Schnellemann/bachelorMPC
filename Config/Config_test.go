package config

import (
	"testing"
)

var path string = "ConfigFiles/configTest.json"

func TestConfigLoad(t *testing.T) {
	expString := "(p1 + p2) + p3"
	ipString := "127.0.1.1"
	portString := "6969"
	ipportParty1 := "127.0.1.1:6969"
	var conf Config = ReadConfig(path)[0]
	if conf.ConstantConfig.Expression != expString {
		t.Errorf("Expression should have been %v, but is %v", expString, conf.ConstantConfig.Expression)
	}
	if conf.ConstantConfig.NumberOfParties != 3 {
		t.Errorf("NumberOfParties should have been %v, but is %v", 3, conf.ConstantConfig.NumberOfParties)
	}
	if conf.ConstantConfig.Ipports[0] != ipportParty1 {
		t.Errorf("Port should have been %v, but is %v", ipportParty1, conf.ConstantConfig.Ipports[0])
	}
	if conf.VariableConfig.Ip != ipString {
		t.Errorf("Ip should have been %v, but is %v", ipString, conf.VariableConfig.Ip)
	}
	if conf.VariableConfig.Port != portString {
		t.Errorf("Port should have been %v, but is %v", portString, conf.VariableConfig.Port)
	}
	if conf.VariableConfig.Secret != 1 {
		t.Errorf("Port should have been %v, but is %v", portString, conf.VariableConfig.Port)
	}

}
