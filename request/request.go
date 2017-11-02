package request

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Interface interface {
	URL() (*url.URL, error)
	Path() string
	Body() string
	Query() Query
	Method() string
	RawURL() string
	Expires() int
	Headers() Headers
	HTTPRequest(string) (*http.Request, error)
	json.Unmarshaler
	json.Marshaler
}

type Request struct {
	url     string
	body    string
	method  string
	expires int
	headers Headers
}

func (r *Request) Path() string {
	return r.parsePathQuery().Path
}

func (r *Request) URL() (*url.URL, error) {
	return url.Parse(r.url)
}

func (r *Request) RawURL() string {
	return r.url
}

func (r *Request) Headers() Headers {
	return r.headers
}

func (r *Request) Method() string {
	return r.method
}

func (r *Request) Body() string {
	return r.body
}

func (r *Request) Expires() int {
	return r.expires
}
