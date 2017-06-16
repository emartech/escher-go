package signer

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	escher "github.com/adamluzsi/escher-go"
)

const defaultExpirationTime = 86400 * time.Second

func (s *signer) SignedURLBy(httpMethod, urlToSign string, expires int) (string, error) {
	uri, err := url.Parse(urlToSign)
	if err != nil {
		return "", err
	}

	date, err := s.config.FormattedDate()
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
	uri.RawQuery = uri.RawQuery + "&" + values.Encode()

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
