package testing

import escher "github.com/adamluzsi/escher-go"

type TestConfigExpected struct {
	Request              escher.Request `json:"request"`
	CanonicalizedRequest string         `json:"canonicalizedRequest"`
	StringToSign         string         `json:"stringToSign"`
	AuthHeader           string         `json:"authHeader"`
	APIKeyID             string         `json:"apiKey"`
	Error                string         `json:"error"`
	URL                  string         `json:"url"`
}
