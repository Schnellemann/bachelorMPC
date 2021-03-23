package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/*
	Config struct for different fields needed to setup network.
*/
type Config struct {
	Expression      string  `json:"expression"`
	NumberOfParties float64 `json:"numberOfParties"`
	Ip              string  `json:"ip"`
	Port            string  `json:"port"`
	PartyNr         float64 `json:"partyNr"`
}

type ProtocolConfig struct {
	Configs []Config `json:"configs"`
}

func readConfig(filepath string) ProtocolConfig {
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
	fmt.Println(conf)

	return conf
}
