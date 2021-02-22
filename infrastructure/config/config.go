package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type (
	// Configuration is struct for holding service's configuration info
	Configuration struct {
		ListenPort string       `json:"ListenPort" validate:"required"`
		Redis      RedisConfig  `json:"RedisConfig" validate:"required"`
		Log        LoggerConfig `json:"Log" validate:"required"`
	}

	// LoggerConfig is a struct for holding logger configuration
	LoggerConfig struct {
		Level       uint32 `json:"Level" validate:"required"`
		ServiceName string `json:"ServiceName" validate:"required"`
		FileName    string `json:"FileName" validate:"required"`
	}

	// RedisConfig is redis configuration
	RedisConfig struct {
		Address  string `json:"Address" validate:"required"`
		Password string `json:"Password" validate:"required"`
		DB       int    `json:"DB" validate:"required"`
	}
)

// New is func for loading app config
func New(configFilePath string) (config Configuration, err error) {
	if _, err = os.Stat("../logs"); os.IsNotExist(err) {
		if err = os.Mkdir("../logs", os.FileMode(0777)); err != nil {
			return
		}
	}

	if config, err = readConfigJSON(configFilePath); err != nil {
		return
	}
	return
}

// readConfigJSON reads config info from JSON file
func readConfigJSON(filePath string) (Configuration, error) {
	log.Printf("Searching JSON config file (%s)", filePath)
	var config Configuration

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Configuration{}, err
	}

	reader := bytes.NewBuffer(contents)
	if err = json.NewDecoder(reader).Decode(&config); err != nil {
		return Configuration{}, fmt.Errorf("error while reading configuration from JSON (%s) error: %w", filePath, err)
	}
	log.Printf("Configuration from JSON (%s) provided\n", filePath)
	return config, nil
}
