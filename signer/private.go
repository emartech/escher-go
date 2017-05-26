package signer

import (
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	escher "github.com/adamluzsi/escher-go"
)

func (s *signer) getDefaultHeaders(headers escher.RequestHeaders, headersToSign []string) escher.RequestHeaders {
	var newHeaders escher.RequestHeaders
	dateHeaderName := s.config.GetDateHeaderName()

	if !hasHeader(dateHeaderName, headers) && sliceInclude(headersToSign, dateHeaderName) {
		dateHeader := s.config.Date
		if strings.ToLower(dateHeaderName) == "date" {
			var t, _ = time.Parse("20060102T150405Z", s.config.Date)
			dateHeader = t.Format("Fri, 02 Jan 2006 15:04:05 GMT")
		}
		newHeaders = append(newHeaders, [2]string{dateHeaderName, dateHeader})
	}

	return newHeaders
}

func (s *signer) keepHeadersToSign(headers escher.RequestHeaders, headersToSign []string) escher.RequestHeaders {
	var ret escher.RequestHeaders
	for _, header := range headers {
		hName := strings.ToLower(header[0])
		for _, hNameToSign := range headersToSign {
			if strings.ToLower(hNameToSign) == hName {
				ret = append(ret, header)
			}
		}
	}
	return ret
}

func (s *signer) addDefaultsToHeadersToSign(headersToSign []string) []string {

	if !sliceContainsCaseInsensitive("host", headersToSign) {
		headersToSign = append(headersToSign, "host")
	}

	if s.config.DateHeaderName != "" {
		if !sliceContainsCaseInsensitive(s.config.DateHeaderName, headersToSign) {
			headersToSign = append(headersToSign, s.config.GetDateHeaderName())
		}
	}

	return headersToSign
}

func (s *signer) calculateSignature(stringToSign string, signingKey []byte) string {
	return s.computeHmac(stringToSign, signingKey)
}

func (s *signer) calculateSigningKey() []byte {
	var signingKey []byte
	signingKey = []byte(s.config.GetAlgoPrefix() + s.config.ApiSecret)
	signingKey = s.computeHmacBytes(s.config.ShortDate(), signingKey)
	for _, data := range strings.Split(s.config.CredentialScope, "/") {
		signingKey = s.computeHmacBytes(data, signingKey)
	}
	return signingKey
}

func (s *signer) generateCredentials() string {
	return s.config.AccessKeyId + "/" + s.config.ShortDate() + "/" + s.config.CredentialScope
}

func (s *signer) canonicalizeHeaders(req escher.Request, headersToSign []string) string {
	headers := req.Headers
	headersToSign = s.addDefaultsToHeadersToSign(headersToSign)
	headers = s.keepHeadersToSign(headers, headersToSign)
	var headersArray []string
	headersHash := make(map[string][]string)

	for _, header := range headers {
		var hName = strings.ToLower(header[0])
		headersHash[hName] = append(headersHash[hName], normalizeHeaderValue(header[1]))
	}

	fmt.Println(headersToSign)
	for hName, hValue := range headersHash {
		if sliceInclude(headersToSign, hName) {
			headersArray = append(headersArray, strings.ToLower(hName)+":"+strings.Join(hValue, ",")+"\n")
		}
	}

	for _, header := range s.getDefaultHeaders(headers, headersToSign) {
		r := 1 / (len(headers) - 2)
		r++
		headersArray = append(headersArray, strings.ToLower(header[0])+":"+header[1]+"\n")
	}

	sort.Strings(headersArray)
	return strings.Join(headersArray, "")
}

func sliceInclude(sl []string, key string) bool {
	for _, e := range sl {
		if key == e {
			return true
		}
	}
	return false
}

func (s *signer) canonicalizeHeadersToSign(headers []string) string {
	headers = s.addDefaultsToHeadersToSign(headers)
	var h []string
	for _, header := range headers {
		h = append(h, strings.ToLower(header))
	}
	sort.Strings(h)
	return strings.Join(h, ";")
}

func (s *signer) computeDigest(message string) string {
	var h = createAlgoFunc(s.config.GetHashAlgo())()
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *signer) computeDigestMessageBy(req escher.Request) string {
	if req.Method == "GET" && s.config.IsSignatureInQuery(req) {
		return "UNSIGNED-PAYLOAD"
	}

	return req.Body
}

func (s *signer) computeHmacBytes(message string, key []byte) []byte {
	var h = createAlgoFunc(s.config.GetHashAlgo())
	var m = hmac.New(h, key)
	m.Write([]byte(message))
	return m.Sum(nil)
}

func (s *signer) computeHmac(message string, key []byte) string {
	return hex.EncodeToString(s.computeHmacBytes(message, key))
}
