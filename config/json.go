package config

import (
	"encoding/json"
)

type jsonMapper struct {
	Date            string `json:"date"`
	HashAlgo        string `json:"hashAlgo"`
	ApiSecret       string `json:"apiSecret"`
	VendorKey       string `json:"vendorKey"`
	AlgoPrefix      string `json:"algoPrefix"`
	AccessKeyId     string `json:"accessKeyId"`
	AuthHeaderName  string `json:"authHeaderName"`
	DateHeaderName  string `json:"dateHeaderName"`
	CredentialScope string `json:"credentialScope"`
}

func ParseJSON(data []byte) (*Config, error) {
	c := &Config{}
	err := mapJSONContentToRequest(c, data)
	return c, err
}

func (c *Config) UnmarshalJSON(data []byte) error {
	return mapJSONContentToRequest(c, data)
}

func mapJSONContentToRequest(c *Config, data []byte) error {
	var j jsonMapper
	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}

	c.Date = j.Date
	c.HashAlgo = j.HashAlgo
	c.ApiSecret = j.ApiSecret
	c.VendorKey = j.VendorKey
	c.AlgoPrefix = j.AlgoPrefix
	c.AccessKeyId = j.AccessKeyId
	c.AuthHeaderName = j.AuthHeaderName
	c.DateHeaderName = j.DateHeaderName
	c.CredentialScope = j.CredentialScope

	return nil
}
