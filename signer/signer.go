package signer

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/EscherAuth/escher"
	"github.com/EscherAuth/escher/request"
)

// Signer is the Escher Signing object interface
type Signer interface {
	SignRequest(request.Interface, []string) (request.Request, error)

	CanonicalizeRequest(request.Interface, []string) string
	GetStringToSign(request.Interface, []string) string
	GenerateHeader(request.Interface, []string) string
	GenerateSignature(request.Interface, []string) string
	SignedURLBy(httpMethod, urlToSign string, expires int) (string, error)
}

type signer struct {
	config escher.Config
}

// New Create a signer object that behaves by the Signer Interface
func New(config escher.Config) Signer {
	return &signer{config}
}

func (s *signer) SignRequest(r request.Interface, headersToSign []string) (request.Request, error) {
	headers := r.Headers()

	var authHeader = s.GenerateHeader(r, headersToSign)
	for _, header := range s.getDefaultHeaders(r) {
		headers = append(headers, header)
	}
	headers = append(headers, [2]string{s.config.AuthHeaderName, authHeader})

	return *request.New(
			r.Method(),
			r.RawURL(),
			headers,
			r.Body(),
			r.Expires()),
		nil
}

func (s *signer) CanonicalizeRequest(r request.Interface, headersToSign []string) string {
	// TODO: remove this shit
	var u = parsePathQuery(r.RawURL())
	parts := make([]string, 0, 6)
	parts = append(parts, r.Method())
	parts = append(parts, canonicalizePath(u.Path))
	parts = append(parts, canonicalizeQuery(u.Query))
	parts = append(parts, s.canonicalizeHeaders(r, headersToSign))
	parts = append(parts, s.canonicalizeHeadersToSign(r, headersToSign))
	parts = append(parts, s.computeDigest(r.Body()))
	canonicalizedRequest := strings.Join(parts, "\n")
	return canonicalizedRequest
}

// TODO: ComposedAlgorithm
func (s *signer) GenerateHeader(r request.Interface, headersToSign []string) string {
	return s.config.AlgoPrefix + "-HMAC-" + s.config.HashAlgo + " " +
		"Credential=" + s.generateCredentials() + ", " +
		"SignedHeaders=" + s.canonicalizeHeadersToSign(r, headersToSign) + ", " +
		"Signature=" + s.GenerateSignature(r, headersToSign)
}

func (s *signer) GenerateSignature(r request.Interface, headersToSign []string) string {
	var stringToSign = s.GetStringToSign(r, headersToSign)
	var signingKey = s.calculateSigningKey()
	return s.calculateSignature(stringToSign, signingKey)
}

func (s *signer) SignedURLBy(httpMethod, urlToSign string, expires int) (string, error) {
	uri, err := url.Parse(urlToSign)

	if err != nil {
		return "", err
	}

	date, err := s.config.DateInEscherFormat()

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

	ereq := request.New(httpMethod, uri.String(), headers, "UNSIGNED-PAYLOAD", expires)

	signature := s.GenerateSignature(ereq, headersToSign)

	values = url.Values{}
	values.Add(s.config.SignatureQueryKey(), signature)
	uri.RawQuery = uri.RawQuery + "&" + values.Encode()

	return uri.String(), nil
}

func (s *signer) GetStringToSign(r request.Interface, headersToSign []string) string {
	return s.config.AlgoPrefix + "-HMAC-" + s.config.HashAlgo + "\n" +
		s.config.Date + "\n" +
		s.config.ShortDate() + "/" + s.config.CredentialScope + "\n" +
		s.computeDigest(s.CanonicalizeRequest(r, headersToSign))
}
