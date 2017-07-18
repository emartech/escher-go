package cases

import "github.com/EscherAuth/escher/request"

type TestConfigExpected struct {
	Request              request.Request `json:"request"`
	CanonicalizedRequest string          `json:"canonicalizedRequest"`
	StringToSign         string          `json:"stringToSign"`
	AuthHeader           string          `json:"authHeader"`
	APIKeyID             string          `json:"apiKey"`
	Error                string          `json:"error"`
	URL                  string          `json:"url"`
}
