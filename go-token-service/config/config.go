package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type appConfig struct {
	PrivateKey      string `json:"privateKey"`
	Network         string `json:"network"`
	InfuraKey       string `json:"infuraKey"`
	ContractAddress string `json:"contractAddress"`
}

var config *appConfig
var once sync.Once
var configFilePath = ""

// SetConfigFilePath allways run once before GetConfig() !
func SetConfigFilePath(filePath string) {
	configFilePath = filePath
}

func GetConfig() *appConfig {
	once.Do(func() {
		config, _ = readConfig()
	})

	return config
}

func readConfig() (*appConfig, error) {
	config := &appConfig{}

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		return config, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, config)
	if err != nil {
		return config, err
	}

	return config, nil
}
