package signer

import (
	"net/url"
	"strconv"
	"strings"

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

func (s *signer) SignRequest(r request.Interface, headersToSign []string) (*request.Request, error) {
	err := s.config.Validate(r)

	if err != nil {
		return nil, err
	}

	headers := r.Headers()

	var authHeader = s.generateHeader(r, headersToSign)
	for _, header := range s.getDefaultHeaders(r) {
		headers = append(headers, header)
	}
	headers = append(headers, [2]string{s.config.AuthHeaderName, authHeader})

	return request.New(
			r.Method(),
			r.RawURL(),
			headers,
			r.Body(),
			r.Expires()),
		nil
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

// TODO add more test to have explicit tests for this not just implicit
func (s *signer) GenerateSignature(r request.Interface, headersToSign []string) string {
	var stringToSign = s.getStringToSign(r, headersToSign)
	var signingKey = s.calculateSigningKey()
	return s.calculateSignature(stringToSign, signingKey)
}
