package escher

import "time"

type Config struct {
	VendorKey       string
	AlgoPrefix      string
	HashAlgo        string
	CredentialScope string
	ApiSecret       string
	AccessKeyId     string
	AuthHeaderName  string
	DateHeaderName  string
	Date            string
}

func (c Config) ShortDate() string {
	return c.Date[:8]
}

func (c Config) Reconfig(t time.Time, apiKeyID, apiSecret string) Config {
	return Config{
		AccessKeyId:     apiKeyID,
		ApiSecret:       apiSecret,
		HashAlgo:        c.HashAlgo,
		VendorKey:       c.VendorKey,
		AlgoPrefix:      c.AlgoPrefix,
		AuthHeaderName:  c.AuthHeaderName,
		DateHeaderName:  c.DateHeaderName,
		CredentialScope: c.CredentialScope,
		Date:            t.Format(time.RFC3339),
	}
}
