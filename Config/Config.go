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

func ReadConfig(filepath string) (configList []*Config) {
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
		configList = append(configList, config)
	}
	return
}

func WriteConfig(filePaths []string, configs []*Config, peersPrComputer int) {

	for i, path := range filePaths {
		sliceOfConfigs := configs[i*peersPrComputer : ((i + 1) * peersPrComputer)]
		file, err := json.MarshalIndent(sliceOfConfigs, "", "")
		if err != nil {
			fmt.Printf("Error in creating")
		}
		err = ioutil.WriteFile(path, file, 0644)
		if err != nil {
			fmt.Printf("Error in writing")
		}
	}

}

func MakeConfigs(ip string, expression string, secrets []int) []*Config {
	var configList []*Config
	listenIpPorts, connectIpPorts := makeIpPortStrings(ip, len(secrets))
	for i, s := range secrets {
		vconfig := VariableConfig{PartyNr: float64(i + 1), Secret: int64(s), ListenIpPort: listenIpPorts[i], ConnectIpPort: connectIpPorts[i]}
		cconfig := ConstantConfig{Expression: expression, NumberOfParties: float64(len(secrets)), Ipports: listenIpPorts}
		config := Config{VariableConfig: vconfig, ConstantConfig: cconfig}
		configList = append(configList, &config)
	}
	return configList
}

func MakeDistributedConfigs(listOfIps []string, nrOfPeers int, secrets []int, expression string) []*Config {
	var configList []*Config
	var connectLists [][]string
	nrOfComputers := len(listOfIps)
	var fullListenList []string
	peersPrComputer := nrOfPeers / nrOfComputers
	if len(secrets) != nrOfPeers {
		fmt.Printf("Length of secrets is %v, and number of peers for this computer is %v", len(secrets), nrOfPeers)
	}
	for i := 0; i < nrOfComputers; i++ {
		computerIp := listOfIps[i]
		listenIpPorts, connectIpPorts := makeIpPortStrings(computerIp, peersPrComputer)
		fullListenList = append(fullListenList, listenIpPorts...)
		connectLists = append(connectLists, connectIpPorts)
	}
	cconfig := ConstantConfig{Expression: expression, NumberOfParties: float64(len(secrets)), Ipports: fullListenList}
	partynr := 0
	for computerNr, connectList := range connectLists {
		for _, connect := range connectList {
			if connect == "" && computerNr != 0 {
				connect = fullListenList[0]
			}
			vconfig := VariableConfig{PartyNr: float64(partynr + 1),
				Secret:        int64(secrets[partynr]),
				ListenIpPort:  fullListenList[partynr],
				ConnectIpPort: connect,
			}
			config := Config{VariableConfig: vconfig, ConstantConfig: cconfig}
			configList = append(configList, &config)
			partynr += 1
		}
	}
	return configList

}

func makeIpPortStrings(ip string, nrOfPeers int) (listenIpPorts []string, connectIpPorts []string) {
	for i := 0; i < nrOfPeers; i++ {
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
	return
}
