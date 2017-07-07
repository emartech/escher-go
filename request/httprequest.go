package request

import (
	"io/ioutil"
	"net/http"
)

func NewFromHTTPRequest(r *http.Request) *Request {

	headers := Headers{}

	for key, values := range r.Header {
		for _, value := range values {
			headers = append(headers, [2]string{key, value})
		}
	}

	body, _ := ioutil.ReadAll(r.Body)

	return New(r.Method, r.URL.String(), headers, string(body), 0)

}
