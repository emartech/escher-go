package request_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/EscherAuth/escher/request"
	"github.com/stretchr/testify/assert"
)

func TestNewFromHTTPRequest(t *testing.T) {

	httpRequest, err := http.NewRequest("GET", "/?k=p", bytes.NewBuffer([]byte("Hello, World!")))

	if err != nil {
		t.Fatal(err)
	}

	httpRequest.Header.Set("X-Testing", "OK")

	escherReqest := request.NewFromHTTPRequest(httpRequest)

	assert.Equal(t, escherReqest.Path(), "/")
	assert.Equal(t, escherReqest.Body(), "Hello, World!")
	assert.Equal(t, escherReqest.Method(), "GET")
	assert.Equal(t, escherReqest.RawURL(), "/?k=p")
	assert.Equal(t, escherReqest.Expires(), 0)
	assert.Equal(t, request.Query{[2]string{"k", "p"}}, escherReqest.Query())
	assert.Equal(t, request.Headers{[2]string{"X-Testing", "OK"}}, escherReqest.Headers())

}
