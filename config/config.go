package config

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/EscherAuth/escher/request"
	"github.com/EscherAuth/escher/utils"
)

type Interface interface {
	json.Unmarshaler
}

type Config struct {
	Date            string
	HashAlgo        string
	ApiSecret       string
	VendorKey       string
	AlgoPrefix      string
	AccessKeyId     string
	AuthHeaderName  string
	DateHeaderName  string
	CredentialScope string
}

func (c Config) ShortDate() string {
	date := c.DateInEscherFormat()

	if date == "" {
		return ""
	}

	return date[:8]
}

func (c Config) Reconfig(date, hashAlgo, credentialScope, apiKeyID, apiSecret string) Config {
	return Config{
		HashAlgo:        hashAlgo,
		AccessKeyId:     apiKeyID,
		ApiSecret:       apiSecret,
		VendorKey:       c.VendorKey,
		AlgoPrefix:      c.AlgoPrefix,
		CredentialScope: credentialScope,
		AuthHeaderName:  c.AuthHeaderName,
		DateHeaderName:  c.DateHeaderName,
		Date:            date,
	}
}

func (c Config) ComposedAlgorithm() string {
	return c.GetAlgoPrefix() + "-HMAC-" + c.GetHashAlgo()
}

func (c Config) DateInEscherFormat() string {
	return c.GetDateWithFormat(utils.EscherDateFormat)
}

func (c Config) DateInHTTPHeaderFormat() string {
	return strings.Replace(c.GetDateWithFormat(utils.HTTPHeaderFormat), "UTC", "GMT", 1)
}

func (c Config) GetDateWithFormat(format string) string {

	if c.Date == "" {
		return time.Now().UTC().Format(format)
	}

	t, err := utils.ParseTime(c.Date)

	if err != nil {
		log.Println("WARNING: malformed date string given, fallback to current time")

		return time.Now().UTC().Format(format)
	}

	return t.UTC().Format(format)

}

func (c Config) GetHashAlgo() string {
	if c.HashAlgo != "" {
		return c.HashAlgo
	}

	return "SHA256"
}

func (c Config) GetAlgoPrefix() string {
	if c.AlgoPrefix != "" {
		return c.AlgoPrefix
	}

	return "ESR"
}

func (c Config) GetVendorKey() string {
	if c.VendorKey != "" {
		return c.VendorKey
	}

	return "Escher"
}

func (c Config) GetAuthHeaderName() string {
	if c.AuthHeaderName != "" {
		return c.AuthHeaderName
	}

	return "X-Escher-Auth"
}

func (c Config) GetDateHeaderName() string {
	if c.DateHeaderName != "" {
		return c.DateHeaderName
	}

	return "X-Escher-Date"
}

func (c Config) IsSigningInQuery(r request.Interface) bool {

	requiredKeys := []string{
		c.QueryKeyFor("Algorithm"),
		c.QueryKeyFor("Credentials"),
		c.QueryKeyFor("Date"),
		c.QueryKeyFor("Expires"),
		c.QueryKeyFor("SignedHeaders"),
	}

	q := r.Query()
	for _, requiredKey := range requiredKeys {
		if !q.IsInclude(requiredKey) {
			return false
		}
	}

	return true
}

func (c Config) SignatureQueryKey() string {
	return c.QueryKeyFor("Signature")
}

func (c Config) QueryKeyFor(key string) string {
	return "X-" + c.GetVendorKey() + "-" + key
}
