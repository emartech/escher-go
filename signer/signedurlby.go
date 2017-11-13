package signer

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/EscherAuth/escher/request"
)

func (s *signer) SignedURLBy(httpMethod, urlToSign string, expires int) (string, error) {
	uri, err := url.Parse(urlToSign)

	if err != nil {
		return "", err
	}

	headers := [][2]string{[2]string{"host", uri.Host}}
	headersToSign := []string{"host"}

	values := url.Values{}
	values.Add(s.config.QueryKeyFor("Algorithm"), s.config.ComposedAlgorithm())
	values.Add(s.config.QueryKeyFor("Credentials"), s.generateCredentials())
	values.Add(s.config.QueryKeyFor("Date"), s.config.DateInEscherFormat())
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
