package signer

import escher "github.com/adamluzsi/escher-go"

type Signer interface {
	CanonicalizeRequest(escher.Request, []string) string
	GetStringToSign(escher.Request, []string) string
	GenerateHeader(escher.Request, []string) string
	SignRequest(escher.Request, []string) escher.Request
	GenerateSignature(escher.Request, []string) string
	SignedURLBy(httpMethod, urlToSign string, expires int) (string, error)
}

type signer struct {
	config escher.Config
}

func New(config escher.Config) Signer {
	return &signer{config}
}
