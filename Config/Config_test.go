package config

import (
	"testing"
)

var path string = "ConfigFiles/configTest.json"

func TestConfigLoad(t *testing.T) {
	expString := "(p1 + p2) + p3"
	listen := "127.0.1.1:6969"
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
	if conf.VariableConfig.ListenIpPort != listen {
		t.Errorf("Ip should have been %v, but is %v", listen, conf.VariableConfig.ListenIpPort)
	}
	if conf.VariableConfig.ConnectIpPort != "" {
		t.Errorf("Port should have been %v, but is %v", "empty string", conf.VariableConfig.ConnectIpPort)
	}
	if conf.VariableConfig.Secret != 1 {
		t.Errorf("Port should have been %v, but is %v", portString, conf.VariableConfig.Secret)
	}

}
