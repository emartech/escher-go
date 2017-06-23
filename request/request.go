package request

import (
	"encoding/json"
	"net/url"
)

type Query [][2]string

type Request interface {
	URL() *url.URL
	Path() string
	Body() string
	Query() Query
	Method() string
	RawURL() string
	Expires() int
	Headers() Headers
	json.Unmarshaler
}

type request struct {
	// UniversalResourceLocator *url.URL

	url     string
	body    string
	method  string
	expires int
	headers Headers
}

func (r *request) Path() string {
	return r.parsePathQuery().Path
}

func (r *request) Query() Query {
	return r.parsePathQuery().Query
}

func (r *request) URL() *url.URL {
	u, _ := url.Parse(r.url)
	return u
}

func (r *request) RawURL() string {
	return r.url
}

func (r *request) Headers() Headers {
	return r.headers
}

func (r *request) Method() string {
	return r.method
}

func (r *request) Body() string {
	return r.body
}

func (r *request) Expires() int {
	return r.expires
}
