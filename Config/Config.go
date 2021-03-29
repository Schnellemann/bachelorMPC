package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	VariableConfig VariableConfig
	ConstantConfig ConstantConfig
}

/*
	Config struct for different fields needed to setup network.
*/
type VariableConfig struct {
	ListenIpPort  string  `json:"listen"`
	ConnectIpPort string  `json:"connect"`
	PartyNr       float64 `json:"partyNr"`
	Secret        int64   `json:"secret"`
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

func MakeConfigs(ip string, expression string, secrets []int) []*Config {
	var configList []*Config
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

		vconfig := VariableConfig{PartyNr: float64(i + 1), Secret: int64(s), ListenIpPort: listenIpPorts[i], ConnectIpPort: connectIpPorts[i]}
		cconfig := ConstantConfig{Expression: expression, NumberOfParties: float64(len(secrets)), Ipports: listenIpPorts}
		config := Config{VariableConfig: vconfig, ConstantConfig: cconfig}
		configList = append(configList, &config)
	}
	return configList
}
