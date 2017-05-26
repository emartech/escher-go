package signer

import escher "github.com/adamluzsi/escher-go"

type Signer interface {
	CanonicalizeRequest(request escher.Request, headersToSign []string) string
	GetStringToSign(request escher.Request, headersToSign []string) string
	GenerateHeader(request escher.Request, headersToSign []string) string
	SignRequest(request escher.Request, headersToSign []string) escher.Request
	GenerateSignature(request escher.Request, headersToSign []string) string
}

type signer struct {
	config escher.Config
}

func New(config escher.Config) Signer {
	return &signer{config}
}
