package signer

import (
	"fmt"
	"os"
	"strings"

	escher "github.com/adamluzsi/escher-go"
)

func (s *signer) SignRequest(request escher.Request, headersToSign []string) escher.Request {
	var authHeader = s.GenerateHeader(request, headersToSign)
	for _, header := range s.getDefaultHeaders(request.Headers, headersToSign) {
		request.Headers = append(request.Headers, header)
	}
	request.Headers = append(request.Headers, [2]string{s.config.GetAuthHeaderName(), authHeader})
	return request
}

var inDebug = os.Getenv("DEBUG") != ""

func (s *signer) CanonicalizeRequest(request escher.Request, headersToSign []string) string {
	var url = parsePathQuery(request.Url)
	parts := make([]string, 0)
	parts = append(parts, request.Method)
	parts = append(parts, canonicalizePath(url.Path))
	parts = append(parts, s.canonicalizeQuery(request))
	parts = append(parts, s.canonicalizeHeaders(request, headersToSign))
	parts = append(parts, s.canonicalizeHeadersToSign(headersToSign))
	parts = append(parts, s.computeDigest(s.computeDigestMessageBy(request)))
	canonicalizedRequest := strings.Join(parts, "\n")

	if inDebug {
		fmt.Println(canonicalizedRequest)
		fmt.Println("--\n\n")
		fmt.Println("--\n\n")
	}

	return canonicalizedRequest
}

func (s *signer) GenerateHeader(request escher.Request, headersToSign []string) string {
	return s.config.GetAlgoPrefix() + "-HMAC-" + s.config.GetHashAlgo() + " " +
		"Credential=" + s.generateCredentials() + ", " +
		"SignedHeaders=" + s.canonicalizeHeadersToSign(headersToSign) + ", " +
		"Signature=" + s.GenerateSignature(request, headersToSign)
}

func (s *signer) GenerateSignature(request escher.Request, headersToSign []string) string {
	var stringToSign = s.GetStringToSign(request, headersToSign)
	var signingKey = s.calculateSigningKey()
	return s.calculateSignature(stringToSign, signingKey)
}
