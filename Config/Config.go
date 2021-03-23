package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	VariableConfig VariableConfig
	ConstantConfig ConstantConfig
}

/*
	Config struct for different fields needed to setup network.
*/
type VariableConfig struct {
	Ip      string  `json:"ip"`
	Port    string  `json:"port"`
	PartyNr float64 `json:"partyNr"`
	Secret  int64   `json:"secret"`
}

type ConstantConfig struct {
	Expression      string   `json:"expression"`
	NumberOfParties float64  `json:"numberOfParties"`
	Ipports         []string `json:"ipports"`
}

type ProtocolConfig struct {
	VariableConfigs []VariableConfig `json:"variableconfigs"`
	ConstantConfig  ConstantConfig   `json:"constantconfig"`
}

func ReadConfig(filepath string) (configList []Config) {
	jsonFile, err := os.Open(filepath)
	defer jsonFile.Close()
	if err != nil {
		fmt.Println("Error opening config file")
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading config file")
	}
	var conf ProtocolConfig
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		fmt.Println("Error unmarshalling json")
		fmt.Println("Error:", err)
	}
	for _, varConf := range conf.VariableConfigs {
		config := new(Config)
		config.VariableConfig = varConf
		config.ConstantConfig = conf.ConstantConfig
		configList = append(configList, *config)
	}
	return
}
