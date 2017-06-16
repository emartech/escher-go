package signer

import (
	"strings"

	escher "github.com/adamluzsi/escher-go"
)

func (s *signer) SignRequest(request escher.Request, headersToSign []string) escher.Request {
	var authHeader = s.GenerateHeader(request, headersToSign)
	for _, header := range s.getDefaultHeaders(request) {
		request.Headers = append(request.Headers, header)
	}
	request.Headers = append(request.Headers, [2]string{s.config.AuthHeaderName, authHeader})
	return request
}

func (s *signer) CanonicalizeRequest(request escher.Request, headersToSign []string) string {
	var url = parsePathQuery(request.Url)
	parts := make([]string, 0, 6)
	parts = append(parts, request.Method)
	parts = append(parts, canonicalizePath(url.Path))
	parts = append(parts, canonicalizeQuery(url.Query))
	parts = append(parts, s.canonicalizeHeaders(request, headersToSign))
	parts = append(parts, s.canonicalizeHeadersToSign(request, headersToSign))
	parts = append(parts, s.computeDigest(request.Body))
	canonicalizedRequest := strings.Join(parts, "\n")
	return canonicalizedRequest
}

// TODO: ComposedAlgorithm
func (s *signer) GenerateHeader(request escher.Request, headersToSign []string) string {
	return s.config.AlgoPrefix + "-HMAC-" + s.config.HashAlgo + " " +
		"Credential=" + s.generateCredentials() + ", " +
		"SignedHeaders=" + s.canonicalizeHeadersToSign(request, headersToSign) + ", " +
		"Signature=" + s.GenerateSignature(request, headersToSign)
}

func (s *signer) GenerateSignature(request escher.Request, headersToSign []string) string {
	var stringToSign = s.GetStringToSign(request, headersToSign)
	var signingKey = s.calculateSigningKey()
	return s.calculateSignature(stringToSign, signingKey)
}
