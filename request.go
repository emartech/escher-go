package escher

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
)

type Request struct {
	Method  string         `json:"method"`
	Url     string         `json:"url"`
	Headers RequestHeaders `json:"headers"`
	Body    string         `json:"body"`
	Expires int            `json:"expires"`
}

type QueryParts [][2]string

func (r Request) Path() string {
	return r.parsePathQuery(r.Url).Path
}

func (r Request) Query() [][2]string {
	return r.parsePathQuery(r.Url).Query
}

func (r Request) QueryParts() (QueryParts, error) {
	u, err := url.Parse(r.Url)

	if err != nil {
		return nil, err
	}

	return transformQueryValues(u.Query()), nil
}

func transformQueryValues(queryValues url.Values) [][2]string {
	queryParts := make([][2]string, 0)

	for key, values := range queryValues {
		for _, value := range values {
			queryParts = append(queryParts, [2]string{key, value})
		}
	}

	return queryParts
}

func (r *Request) DelQueryValueByKey(key string) error {
	uri, err := url.Parse(r.Url)

	if err != nil {
		return err
	}

	values := uri.Query()
	values.Del(key)
	uri.RawQuery = values.Encode()

	r.Url = uri.String()

	return nil
}

// ToDo: remove this later

func (r Request) parsePathQuery(pathAndQuery string) parsedPathQuery {
	u, err := url.Parse(pathAndQuery)
	var p parsedPathQuery
	s := strings.SplitN(pathAndQuery, "?", 2)
	p.Path = s[0]

	if err == nil {
		p.Query = transformQueryValues(u.Query())
	} else {

		p.Query = parseQuery(s[1])

	}

	if err == nil && u.Scheme != "" {
		p.Path = u.Path
	}

	return p
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

func hasHeader(headerName string, headers RequestHeaders) bool {
	headerName = strings.ToLower(headerName)
	for _, header := range headers {
		hName := strings.ToLower(header[0])
		if hName == headerName {
			return true
		}
	}
	return false
}

// ToDo: remove this later
type parsedPathQuery struct {
	Path  string
	Query [][2]string
}

func parseQuery(query string) [][2]string {
	var q [][2]string
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
