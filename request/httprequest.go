package request

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewFromHTTPRequest(r *http.Request) (*Request, error) {

	body, err := bodyStringFrom(r)

	if err != nil {
		return nil, err
	}

	return New(r.Method, r.URL.RequestURI(), createEscherHeadersFromHTTPHeaders(r), body, 36000), nil

}

func (r *Request) HTTPRequest(baseURL string) (*http.Request, error) {

	u, err := url.Parse(baseURL)

	if err != nil {
		return nil, err
	}

	rURL, err := r.URL()

	if err != nil {
		return nil, err
	}

	mergeURLPath(rURL, u)

	httpRequest, err := http.NewRequest(r.method, u.String(), bytes.NewBuffer([]byte(r.body)))

	if err != nil {
		return httpRequest, err
	}

	for _, keyValuePair := range r.headers {
		key := keyValuePair[0]
		value := keyValuePair[1]

		httpRequest.Header.Add(key, value)
	}

	return httpRequest, nil
}

func bodyStringFrom(r *http.Request) (string, error) {

	if r.Body == nil {
		return "", nil
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return string(body), nil

}

func createEscherHeadersFromHTTPHeaders(r *http.Request) Headers {
	headers := Headers{}

	for key, values := range r.Header {
		for _, value := range values {
			headers = append(headers, [2]string{key, value})
		}
	}

	if r.Header.Get("host") == "" {
		// https://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.23
		headers = append(headers, [2]string{"host", r.URL.Host})
	}

	return headers
}
