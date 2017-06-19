package signer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/EscherAuth/escher"
)

// Signer is the Escher Signing object interface
type Signer interface {
	// CanonicalizeRequest creates a unified representing form from the request in a string
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

// New Create a signer object that behaves by the Signer Interface
func New(config escher.Config) Signer {
	return &signer{config}
}

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
	fmt.Println(s.canonicalizeHeadersToSign(request, headersToSign))
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

func (s *signer) SignedURLBy(httpMethod, urlToSign string, expires int) (string, error) {
	uri, err := url.Parse(urlToSign)

	if err != nil {
		return "", err
	}

	date, err := s.config.GetDate()

	if err != nil {
		return "", err
	}

	headers := [][2]string{[2]string{"host", uri.Host}}
	headersToSign := []string{"host"}

	values := url.Values{}
	values.Add(s.config.QueryKeyFor("Algorithm"), s.config.ComposedAlgorithm())
	values.Add(s.config.QueryKeyFor("Credentials"), s.generateCredentials())
	values.Add(s.config.QueryKeyFor("Date"), date)
	values.Add(s.config.QueryKeyFor("Expires"), strconv.Itoa(expires))
	values.Add(s.config.QueryKeyFor("SignedHeaders"), strings.Join(headersToSign, ";"))

	if uri.RawQuery == "" {
		uri.RawQuery = values.Encode()
	} else {
		uri.RawQuery = uri.RawQuery + "&" + values.Encode()
	}

	fmt.Println("| | |")
	fmt.Println(uri.String())
	ereq := escher.Request{
		Method:  httpMethod,
		Url:     uri.String(),
		Headers: headers,
		Body:    "UNSIGNED-PAYLOAD",
		Expires: expires,
	}

	signature := s.GenerateSignature(ereq, headersToSign)

	values = url.Values{}
	values.Add(s.config.SignatureQueryKey(), signature)
	uri.RawQuery = uri.RawQuery + "&" + values.Encode()

	return uri.String(), nil
}
