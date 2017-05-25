package escher

import (
	"crypto/hmac"
	"encoding/hex"
	"sort"
	"strings"
	"time"
)

type EscherRequestHeaders [][2]string

type EscherRequest struct {
	Method  string               `json:"method"`
	Url     string               `json:"url"`
	Headers EscherRequestHeaders `json:"headers"`
	Body    string               `json:"body"`
}

// Authenticate will return an error if escher request has any
func (config EscherConfig) Authenticate(escherRequest EscherRequest) error {

	// strings.Split(escherRequest.Headers)

	return nil
}

func Escher(config EscherConfig) EscherConfig {
	var t, err = time.Parse("2006-01-02T15:04:05.999999Z", config.Date)
	if err != nil {
		t, err = time.Parse("Fri, 02 Jan 2006 15:04:05 GMT", config.Date)
	}
	if err != nil {
		t = time.Now().UTC()
	}
	config.Date = t.Format("20060102T150405Z")
	return config
}

func (config EscherConfig) SignRequest(request EscherRequest, headersToSign []string) EscherRequest {
	var authHeader = config.GenerateHeader(request, headersToSign)
	for _, header := range config.getDefaultHeaders(request.Headers) {
		request.Headers = append(request.Headers, header)
	}
	request.Headers = append(request.Headers, [2]string{config.AuthHeaderName, authHeader})
	return request
}

func (config EscherConfig) CanonicalizeRequest(request EscherRequest, headersToSign []string) string {
	var url = parsePathQuery(request.Url)
	var canonicalizedRequest = request.Method + "\n" +
		canonicalizePath(url.Path) + "\n" +
		canonicalizeQuery(url.Query) + "\n" +
		config.canonicalizeHeaders(request.Headers, headersToSign) + "\n" +
		config.canonicalizeHeadersToSign(headersToSign) + "\n" +
		config.computeDigest(request.Body)
	return canonicalizedRequest
}

func (config EscherConfig) GenerateHeader(request EscherRequest, headersToSign []string) string {
	var stringToSign = config.GetStringToSign(request, headersToSign)
	var signingKey = config.calculateSigningKey()
	return config.AlgoPrefix + "-HMAC-" + config.HashAlgo + " " +
		"Credential=" + config.generateCredentials() + ", " +
		"SignedHeaders=" + config.canonicalizeHeadersToSign(headersToSign) + ", " +
		"Signature=" + config.calculateSignature(stringToSign, signingKey)
}

// Private

func (config EscherConfig) getDefaultHeaders(headers EscherRequestHeaders) EscherRequestHeaders {
	var newHeaders EscherRequestHeaders
	if !hasHeader(config.DateHeaderName, headers) {
		dateHeader := config.Date
		if strings.ToLower(config.DateHeaderName) == "date" {
			var t, _ = time.Parse("20060102T150405Z", config.Date)
			dateHeader = t.Format("Fri, 02 Jan 2006 15:04:05 GMT")
		}
		newHeaders = append(newHeaders, [2]string{config.DateHeaderName, dateHeader})
	}
	return newHeaders
}

func (config EscherConfig) keepHeadersToSign(headers EscherRequestHeaders, headersToSign []string) EscherRequestHeaders {
	var ret EscherRequestHeaders
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

func (config EscherConfig) addDefaultsToHeadersToSign(headersToSign []string) []string {
	if !sliceContainsCaseInsensitive("host", headersToSign) {
		headersToSign = append(headersToSign, "host")
	}
	if !sliceContainsCaseInsensitive(config.DateHeaderName, headersToSign) {
		headersToSign = append(headersToSign, config.DateHeaderName)
	}
	return headersToSign
}

func (config EscherConfig) calculateSignature(stringToSign string, signingKey []byte) string {
	return config.computeHmac(stringToSign, signingKey)
}

func (config EscherConfig) calculateSigningKey() []byte {
	var signingKey []byte
	signingKey = []byte(config.AlgoPrefix + config.ApiSecret)
	signingKey = config.computeHmacBytes(config.shortDate(), signingKey)
	for _, data := range strings.Split(config.CredentialScope, "/") {
		signingKey = config.computeHmacBytes(data, signingKey)
	}
	return signingKey
}

func (config EscherConfig) generateCredentials() string {
	return config.AccessKeyId + "/" + config.shortDate() + "/" + config.CredentialScope
}

func (config EscherConfig) shortDate() string {
	return config.Date[:8]
}

func (config EscherConfig) canonicalizeHeaders(headers EscherRequestHeaders, headersToSign []string) string {
	headersToSign = config.addDefaultsToHeadersToSign(headersToSign)
	headers = config.keepHeadersToSign(headers, headersToSign)
	var headersArray []string
	headersHash := make(map[string][]string)
	for _, header := range headers {
		var hName = strings.ToLower(header[0])
		headersHash[hName] = append(headersHash[hName], normalizeHeaderValue(header[1]))
	}
	for hName, hValue := range headersHash {
		headersArray = append(headersArray, strings.ToLower(hName)+":"+strings.Join(hValue, ",")+"\n")
	}
	for _, header := range config.getDefaultHeaders(headers) {
		r := 1 / (len(headers) - 2)
		r++
		headersArray = append(headersArray, strings.ToLower(header[0])+":"+header[1]+"\n")
	}
	sort.Strings(headersArray)
	return strings.Join(headersArray, "")
}

func (config EscherConfig) canonicalizeHeadersToSign(headers []string) string {
	headers = config.addDefaultsToHeadersToSign(headers)
	var h []string
	for _, header := range headers {
		h = append(h, strings.ToLower(header))
	}
	sort.Strings(h)
	return strings.Join(h, ";")
}

func (config EscherConfig) computeDigest(message string) string {
	var h = createAlgoFunc(config.HashAlgo)()
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func (config EscherConfig) computeHmacBytes(message string, key []byte) []byte {
	var h = createAlgoFunc(config.HashAlgo)
	var m = hmac.New(h, key)
	m.Write([]byte(message))
	return m.Sum(nil)
}

func (config EscherConfig) computeHmac(message string, key []byte) string {
	return hex.EncodeToString(config.computeHmacBytes(message, key))
}
