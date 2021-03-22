package config

import (
	"testing"
)

var path string = "configTest.json"

func TestConfigLoad(t *testing.T) {
	expString := "(p1 + p2) + p3"
	var conf Config = readConfig(path)
	if conf.Expression != expString {
		t.Errorf("String should have been %v, but is %v", expString, conf.Expression)
	}
}
