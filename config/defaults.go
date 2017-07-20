package config

const (
	defaultAlgoPrefix     = "ESR"
	defaultHashAlgo       = "SHA256"
	defaultVendorKey      = "Escher"
	defaultAuthHeaderName = "X-Escher-Auth"
	defaultDateHeaderName = "X-Escher-Date"
)

func setDefaults(c *Config) {

	if c.AlgoPrefix == "" {
		c.AlgoPrefix = defaultAlgoPrefix
	}

	if c.HashAlgo == "" {
		c.HashAlgo = defaultHashAlgo
	}

	if c.VendorKey == "" {
		c.VendorKey = defaultVendorKey
	}

	if c.AuthHeaderName == "" {
		c.AuthHeaderName = defaultAuthHeaderName
	}

	if c.DateHeaderName == "" {
		c.DateHeaderName = defaultDateHeaderName
	}

}
