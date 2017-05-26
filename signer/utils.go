package signer

import (
	"crypto/sha256"
	"crypto/sha512"
	"strings"

	"github.com/PuerkitoBio/purell"
	escher "github.com/adamluzsi/escher-go"

	"hash"
	"net/url"
	"regexp"
)

type parsedPathQuery struct {
	Path  string
	Query requestQuery
}

type requestQuery [][2]string

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

func (s *signer) canonicalizeQuery(r escher.Request) string {
	u, err := url.Parse(r.Url)
	if err != nil {
		return ""
	}

	isSignatureInQuery := s.config.IsSignatureInQuery(r)
	signatureQueryKey := s.config.SignatureQueryKey()

	queryValues := make(url.Values)
	for key, values := range u.Query() {
		if isSignatureInQuery && key == signatureQueryKey {
			continue
		}
		queryValues[key] = values
	}

	return queryValues.Encode()
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

func parsePathQuery(pathAndQuery string) parsedPathQuery {
	var p parsedPathQuery
	s := strings.SplitN(pathAndQuery, "?", 2)
	p.Path = s[0]
	if len(s) > 1 {
		p.Query = parseQuery(s[1])
	}
	return p
}

func (s *signer) GetStringToSign(request escher.Request, headersToSign []string) string {
	return s.config.GetAlgoPrefix() + "-HMAC-" + s.config.GetHashAlgo() + "\n" +
		s.config.Date + "\n" +
		s.config.ShortDate() + "/" + s.config.CredentialScope + "\n" +
		s.computeDigest(s.CanonicalizeRequest(request, headersToSign))
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

func parseQuery(query string) requestQuery {
	var q requestQuery
	for _, param := range strings.Split(query, "&") {
		var kv = strings.SplitN(param, "=", 2)
		var kv2 [2]string
		kv2[0] = queryUnescape(kv[0])
		if len(kv) > 1 {
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
			ret = append(ret, unhex(s[i+1])<<4|unhex(s[i+2]))
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
	path = queryUnescape(purell.NormalizeURL(&u, purell.FlagRemoveDotSegments|purell.FlagRemoveDuplicateSlashes))
	return path
}

func hasHeader(headerName string, headers escher.RequestHeaders) bool {
	headerName = strings.ToLower(headerName)
	for _, header := range headers {
		hName := strings.ToLower(header[0])
		if hName == headerName {
			return true
		}
	}
	return false
}
