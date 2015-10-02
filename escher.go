package escher

import (
  "hash"
  "crypto/hmac"
  "crypto/sha256"
  "crypto/sha512"
  "encoding/hex"
  "strings"
  "sort"
  "time"
  "net/url"
  "regexp"
  . "github.com/PuerkitoBio/purell"
)

type parsedPathQuery struct {
  Path string
  Query EshcerRequestQuery
}

type EscherConfig struct {
  VendorKey string
  AlgoPrefix string
  HashAlgo string
  CredentialScope string
  ApiSecret string
  AccessKeyId string
  AuthHeaderName string
  DateHeaderName string
  Date string
}

type EscherRequest struct {
  Method string
  Url string
  Headers EscherRequestHeaders
  Body string
}

type EscherRequestHeaders [][2]string
type EshcerRequestQuery [][2]string

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
  request.Headers = append(request.Headers, [2]string { config.AuthHeaderName, authHeader })
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
    headersArray = append(headersArray, strings.ToLower(hName) + ":" + strings.Join(hValue, ",") + "\n")
  }
  for _, header := range config.getDefaultHeaders(headers) {
    r := 1/(len(headers) - 2)
    r++
    headersArray = append(headersArray, strings.ToLower(header[0]) + ":" + header[1] + "\n")
  }
  sort.Strings(headersArray)
  return strings.Join(headersArray, "")
}

func normalizeHeaderValue(value string) string {
  var valueArray []string
  var betweenQuotes bool = false
  var reWhiteSpace = regexp.MustCompile(" +")
  for _, part := range strings.Split(value, "\"") {
    if !betweenQuotes {
      part = reWhiteSpace.ReplaceAllString(part, " ")
    }
    valueArray = append(valueArray, part)
    betweenQuotes = !betweenQuotes
  }
  return strings.TrimSpace(strings.Join(valueArray, "\""))
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


func canonicalizeQuery(query EshcerRequestQuery) string {
  var q []string
  for _, kv := range query {
    q = append(q, strings.Replace(url.QueryEscape(kv[0]), "+", "%20", -1) + "=" + url.QueryEscape(kv[1]))
  }
  sort.Strings(q)
  return strings.Join(q, "&")
}

func createAlgoFunc(hashAlgo string) func() hash.Hash {
  var h func() hash.Hash
  if hashAlgo == "SHA256" {
    h = sha256.New
  }
  if hashAlgo == "SHA512" {
    h = sha512.New
  }
  return h
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

func parsePathQuery(pathAndQuery string) parsedPathQuery {
  var p parsedPathQuery
  s := strings.SplitN(pathAndQuery, "?", 2)
  p.Path = s[0]
  if len(s) > 1 {
    p.Query = parseQuery(s[1])
  }
  return p
}

func (config EscherConfig) GetStringToSign(request EscherRequest, headersToSign []string) string {
  return config.AlgoPrefix + "-HMAC-" + config.HashAlgo + "\n" +
    config.Date + "\n" +
    config.shortDate() + "/" + config.CredentialScope + "\n" +
    config.computeDigest(config.CanonicalizeRequest(request, headersToSign))
}

func (config EscherConfig) GenerateHeader(request EscherRequest, headersToSign []string) string {
  var stringToSign = config.GetStringToSign(request, headersToSign)
  var signingKey = config.calculateSigningKey()
  return config.AlgoPrefix + "-HMAC-" + config.HashAlgo + " " +
    "Credential=" + config.generateCredentials() + ", " +
    "SignedHeaders=" + config.canonicalizeHeadersToSign(headersToSign) + ", " +
    "Signature=" + config.calculateSignature(stringToSign, signingKey)
}

func sliceContainsCaseInsensitive(needle string, stack []string) bool {
  needle = strings.ToLower(needle)
  for _, item := range stack {
    if strings.ToLower(item) == needle {
      return true
    }
  }
  return false
}

func (config EscherConfig) getDefaultHeaders(headers EscherRequestHeaders) EscherRequestHeaders {
  var newHeaders EscherRequestHeaders
  if !hasHeader(config.DateHeaderName, headers) {
    dateHeader := config.Date
    if strings.ToLower(config.DateHeaderName) == "date" {
      var t, _ = time.Parse("20060102T150405Z", config.Date)
      dateHeader = t.Format("Fri, 02 Jan 2006 15:04:05 GMT")
    }
    newHeaders = append(newHeaders, [2]string { config.DateHeaderName, dateHeader })
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

func parseQuery(query string) EshcerRequestQuery {
  var q EshcerRequestQuery
  for _, param := range strings.Split(query, "&") {
    var kv = strings.SplitN(param, "=", 2)
    var kv2 [2]string
    kv2[0] = queryUnescape(kv[0])
    if (len(kv) > 1) {
      kv2[1] = queryUnescape(kv[1])
    }
    q = append(q, kv2)
  }
  return q
}

func queryUnescape(s string) string {
  var ret []byte
  for i := 0; i < len(s); {
    if s[i] == '%' && i+2 < len(s) && ishex(s[i+1]) && ishex(s[i+2]) {
      ret = append(ret, unhex(s[i+1])<<4 | unhex(s[i+2]))
      i += 2
    } else if s[i] == '+' {
      ret = append(ret, ' ')
    } else {
      ret = append(ret, s[i])
    }
    i++
  }
  return string(ret)
}

func unhex(c byte) byte {
  switch {
    case '0' <= c && c <= '9':
      return c - '0'
    case 'a' <= c && c <= 'f':
      return c - 'a' + 10
    case 'A' <= c && c <= 'F':
      return c - 'A' + 10
  }
  return 0
}

func ishex(c byte) bool {
  return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func canonicalizePath(path string) string {
  var u url.URL
  u.Path = path
  path = queryUnescape(NormalizeURL(&u, FlagRemoveDotSegments | FlagRemoveDuplicateSlashes))
 	return path
}

func hasHeader(headerName string, headers EscherRequestHeaders) bool {
  headerName = strings.ToLower(headerName)
  for _, header := range headers {
    hName := strings.ToLower(header[0])
    if hName == headerName {
      return true
    }
  }
  return false
}
