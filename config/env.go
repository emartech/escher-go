package config

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	escherConfigEnv = "ESCHER_CONFIG"
)

func NewFromENV() (Config, error) {
	c := Config{}

	err := setByConfigJSON(&c)

	if err != nil {
		return c, err
	}

	setValuesFromDifferentEnv(&c)

	if c.CredentialScope == "" {
		return c, errors.New("Credential Scope is missing")
	}

	SetDefaults(&c)

	return c, nil
}

func setByConfigJSON(c *Config) error {
	var err error

	data, configJSONStringIsPresent := os.LookupEnv(escherConfigEnv)

	if configJSONStringIsPresent {
		err = json.Unmarshal([]byte(data), c)
	}

	return err
}

func setValuesFromDifferentEnv(c *Config) {

	algoPrefix, isGiven := os.LookupEnv("ESCHER_ALGO_PREFIX")
	if isGiven && c.AlgoPrefix == "" {
		c.AlgoPrefix = algoPrefix
	}

	hashAlgo, isGiven := os.LookupEnv("ESCHER_HASH_ALGO")
	if isGiven && c.HashAlgo == "" {
		c.HashAlgo = hashAlgo
	}

	vendorKey, isGiven := os.LookupEnv("ESCHER_VENDOR_KEY")
	if isGiven && c.VendorKey == "" {
		c.VendorKey = vendorKey
	}

	authHeaderName, isGiven := os.LookupEnv("ESCHER_AUTH_HEADER_NAME")
	if isGiven && c.AuthHeaderName == "" {
		c.AuthHeaderName = authHeaderName
	}

	dateHeaderName, isGiven := os.LookupEnv("ESCHER_DATE_HEADER_NAME")
	if isGiven && c.DateHeaderName == "" {
		c.DateHeaderName = dateHeaderName
	}

	credentialScope, isGiven := os.LookupEnv("ESCHER_CREDENTIAL_SCOPE")
	if isGiven && c.CredentialScope == "" {
		c.CredentialScope = credentialScope
	}

}
