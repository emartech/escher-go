package signer

import (
	"crypto/hmac"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/EscherAuth/escher/request"
)

func (s *signer) getDefaultHeaders(r request.Interface) request.Headers {
	headers := r.Headers()
	var newHeaders request.Headers

	// TODO: Should I remove date from the headers to sign in IsSigningInQuery case ?
	if !hasHeader(s.config.DateHeaderName, headers) && !s.config.IsSigningInQuery(r) {
		var dateHeader string
		if strings.ToLower(s.config.DateHeaderName) == "date" {
			dateHeader = s.config.DateInHTTPHeaderFormat()
		} else {
			dateHeader = s.config.DateInEscherFormat()
		}

		newHeaders = append(newHeaders, [2]string{s.config.DateHeaderName, dateHeader})
	}

	return newHeaders
}

func (s *signer) keepHeadersToSign(headers request.Headers, headersToSign []string) request.Headers {
	ret := request.Headers{}
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

func (s *signer) addDefaultsToHeadersToSign(r request.Interface, headersToSign []string) []string {

	if !sliceContainsCaseInsensitive("host", headersToSign) {
		headersToSign = append(headersToSign, "host")
	}

	if !s.config.IsSigningInQuery(r) && !sliceContainsCaseInsensitive(s.config.DateHeaderName, headersToSign) {
		headersToSign = append(headersToSign, s.config.DateHeaderName)
	}

	return headersToSign

}

func (s *signer) calculateSignature(stringToSign string, signingKey []byte) string {
	return s.computeHmac(stringToSign, signingKey)
}

func (s *signer) calculateSigningKey() []byte {
	var signingKey []byte
	signingKey = []byte(s.config.AlgoPrefix + s.config.ApiSecret)
	signingKey = s.computeHmacBytes(s.config.ShortDate(), signingKey)
	for _, data := range strings.Split(s.config.CredentialScope, "/") {
		signingKey = s.computeHmacBytes(data, signingKey)
	}
	return signingKey
}

func (s *signer) generateCredentials() string {
	return s.config.AccessKeyId + "/" + s.config.ShortDate() + "/" + s.config.CredentialScope
}

func (s *signer) canonicalizeHeaders(r request.Interface, headersToSign []string) string {
	headers := r.Headers()
	headersToSign = s.addDefaultsToHeadersToSign(r, headersToSign)
	headers = s.keepHeadersToSign(headers, headersToSign)

	var headersArray []string
	headersHash := make(map[string][]string)

	for _, header := range headers {
		var hName = strings.ToLower(header[0])
		headersHash[hName] = append(headersHash[hName], normalizeHeaderValue(header[1]))
	}

	for hName, hValue := range headersHash {
		headersArray = append(headersArray, strings.ToLower(hName)+":"+strings.Join(hValue, ",")+"\n")
	}

	for _, header := range s.getDefaultHeaders(r) {
		headersArray = append(headersArray, strings.ToLower(header[0])+":"+header[1]+"\n")
	}

	sort.Strings(headersArray)
	return strings.Join(headersArray, "")
}

func (s *signer) canonicalizeHeadersToSign(r request.Interface, headersToSign []string) string {
	headers := s.addDefaultsToHeadersToSign(r, headersToSign)

	var loweredHeaders []string
	for _, header := range headers {
		loweredHeaders = append(loweredHeaders, strings.ToLower(header))
	}
	sort.Strings(loweredHeaders)

	return strings.Join(loweredHeaders, ";")
}

func (s *signer) computeDigest(message string) string {
	var h = createAlgoFunc(s.config.HashAlgo)()
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *signer) computeHmacBytes(message string, key []byte) []byte {
	var h = createAlgoFunc(s.config.HashAlgo)
	var m = hmac.New(h, key)
	m.Write([]byte(message))
	return m.Sum(nil)
}

func (s *signer) computeHmac(message string, key []byte) string {
	return hex.EncodeToString(s.computeHmacBytes(message, key))
}

// TODO: ComposedAlgorithm
func (s *signer) generateHeader(r request.Interface, headersToSign []string) string {
	return s.config.AlgoPrefix + "-HMAC-" + s.config.HashAlgo + " " +
		"Credential=" + s.generateCredentials() + ", " +
		"SignedHeaders=" + s.canonicalizeHeadersToSign(r, headersToSign) + ", " +
		"Signature=" + s.GenerateSignature(r, headersToSign)
}

func (s *signer) getStringToSign(r request.Interface, headersToSign []string) string {
	return s.config.AlgoPrefix + "-HMAC-" + s.config.HashAlgo + "\n" +
		s.config.DateInEscherFormat() + "\n" +
		s.config.ShortDate() + "/" + s.config.CredentialScope + "\n" +
		s.computeDigest(s.CanonicalizeRequest(r, headersToSign))
}
