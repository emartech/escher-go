package signer

import (
	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/request"
)

// Signer is the Escher Signing object interface
type Signer interface {
	SignRequest(r request.Interface, headersToSign []string) (*request.Request, error)
	SignedURLBy(httpMethod, urlToSign string, expires int) (string, error)
	GenerateSignature(r request.Interface, headersToSign []string) string
}

type signer struct {
	config config.Config
}

// New Create a signer object that behaves by the Signer Interface
func New(c config.Config) Signer {
	return &signer{c}
}
