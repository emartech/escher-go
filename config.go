package escher

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
