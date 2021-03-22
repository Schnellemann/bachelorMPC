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
	Expression string `json:"expression"`
}

func readConfig(filepath string) Config {
	jsonFile, err := os.Open(filepath)
	defer jsonFile.Close()
	if err != nil {
		fmt.Println("Error opening config file")
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading config file")
	}
	var conf Config
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		fmt.Println("Error unmarshalling json")
		fmt.Println("Error:", err)
	}

	return conf
}
