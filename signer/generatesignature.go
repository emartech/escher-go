package signer

import "github.com/EscherAuth/escher/request"

// TODO add more test to have explicit tests for this not just implicit
func (s *signer) GenerateSignature(r request.Interface, headersToSign []string) string {
	var stringToSign = s.getStringToSign(r, headersToSign)
	var signingKey = s.calculateSigningKey()
	return s.calculateSignature(stringToSign, signingKey)
}
