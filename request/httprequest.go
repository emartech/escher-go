package request

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func NewFromHTTPRequest(r *http.Request) (*Request, error) {

	headers := Headers{}

	for key, values := range r.Header {
		for _, value := range values {
			headers = append(headers, [2]string{key, value})
		}
	}

	// GetBody
	bodyIO, err := r.GetBody()
	if err != nil {
		return nil, err
	}
	defer bodyIO.Close()

	bodyContent, err := ioutil.ReadAll(bodyIO)

	if err != nil {
		return nil, err
	}

	return New(r.Method, r.URL.String(), headers, string(bodyContent), 0), nil

}

func (r *Request) HTTPRequest() (*http.Request, error) {
	bodyIO := bytes.NewBuffer([]byte(r.body))
	httpRequest, err := http.NewRequest(r.method, r.url, bodyIO)

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
