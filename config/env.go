package config

import (
	"encoding/json"
	"os"
)

const (
	escherConfigEnv = "ESCHER_CONFIG"
)

func NewFromENV() (Config, error) {
	config := Config{}

	jsonString, configJSONStringIsPresent := os.LookupEnv(escherConfigEnv)

	if configJSONStringIsPresent {
		json.Unmarshal([]byte(jsonString), &config)
	}

	setDefaults(&config)

	return config, nil
}
