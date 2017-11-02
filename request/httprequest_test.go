package request_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/EscherAuth/escher/request"
	"github.com/stretchr/testify/assert"
)

func TestNewFromHTTPRequest(t *testing.T) {
	t.Parallel()

	httpRequest, err := http.NewRequest("GET", "/?k=p", bytes.NewBuffer([]byte("Hello, World!")))

	if err != nil {
		t.Fatal(err)
	}

	httpRequest.Header.Set("X-Testing", "OK")

	escherReqest, err := request.NewFromHTTPRequest(httpRequest)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, escherReqest.Path(), "/")
	assert.Equal(t, escherReqest.Body(), "Hello, World!")
	assert.Equal(t, escherReqest.Method(), "GET")
	assert.Equal(t, escherReqest.RawURL(), "/?k=p")
	assert.Equal(t, escherReqest.Expires(), 0)
	assert.Equal(t, request.Query{[2]string{"k", "p"}}, escherReqest.Query())
	assert.Equal(t, request.Headers{[2]string{"X-Testing", "OK"}}, escherReqest.Headers())

}

func TestNewFromHTTPRequest_HTTPRequestIncludesSchemaAndOtherImportantParameters_OnlyPathIsUsed(t *testing.T) {
	t.Parallel()

	httpRequest, err := http.NewRequest("GET", "https://example.org/?k=p", nil)

	if err != nil {
		t.Fatal(err)
	}

	escherReqest, err := request.NewFromHTTPRequest(httpRequest)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, escherReqest.Path(), "/")
	assert.Equal(t, escherReqest.RawURL(), "/?k=p")
	assert.Equal(t, request.Query{[2]string{"k", "p"}}, escherReqest.Query())

}

func TestNewFromHTTPRequest_TheRequestBodyIsNil_EmptyStringUsed(t *testing.T) {
	t.Parallel()

	httpRequest, err := http.NewRequest("GET", "/?k=p", nil)

	if err != nil {
		t.Fatal(err)
	}

	escherReqest, err := request.NewFromHTTPRequest(httpRequest)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, escherReqest.Path(), "/")
	assert.Equal(t, escherReqest.Body(), "")
	assert.Equal(t, escherReqest.Method(), "GET")
	assert.Equal(t, escherReqest.RawURL(), "/?k=p")
	assert.Equal(t, escherReqest.Expires(), 0)
	assert.Equal(t, request.Query{[2]string{"k", "p"}}, escherReqest.Query())
	assert.Equal(t, request.Headers{}, escherReqest.Headers())

}

func TestNewFromHTTPRequest_EscherRequestMade_HTTPBodyStillContainsValueLikeItIsUnTouched(t *testing.T) {
	t.Parallel()

	expectedBodyString := "Hello, World!"
	httpRequest, err := http.NewRequest("GET", "/?k=p", bytes.NewBuffer([]byte(expectedBodyString)))

	if err != nil {
		t.Fatal(err)
	}

	httpRequest.Header.Set("X-Testing", "OK")

	escherReqest, err := request.NewFromHTTPRequest(httpRequest)

	if err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadAll(httpRequest.Body)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(content), expectedBodyString)
	assert.Equal(t, escherReqest.Body(), expectedBodyString)
	assert.Equal(t, string(content), escherReqest.Body())

}

func TestHTTPRequest(t *testing.T) {
	t.Parallel()

	newBodyIO := func() *bytes.Buffer { return bytes.NewBuffer([]byte("Hello you awesome!")) }
	originalHTTPRequest, err := http.NewRequest("GET", "/?k=p", newBodyIO())

	if err != nil {
		t.Fatal(err)
	}

	originalHTTPRequest.Header.Set("X-Testing", "OK")

	escherReqest, err := request.NewFromHTTPRequest(originalHTTPRequest)

	if err != nil {
		t.Fatal(err)
	}

	actuallyHTTPRequest, err := escherReqest.HTTPRequest("http://www.example.com")

	if err != nil {
		t.Fatal(err)
	}

	expectedHTTPRequest, err := http.NewRequest("GET", "http://www.example.com/?k=p", newBodyIO())

	if err != nil {
		t.Fatal(err)
	}

	expectedHTTPRequest.Header.Set("X-Testing", "OK")

	assert.Equal(t, expectedHTTPRequest.Method, actuallyHTTPRequest.Method)
	assert.Equal(t, expectedHTTPRequest.URL, actuallyHTTPRequest.URL)
	assert.Equal(t, expectedHTTPRequest.Proto, actuallyHTTPRequest.Proto)
	assert.Equal(t, expectedHTTPRequest.ProtoMajor, actuallyHTTPRequest.ProtoMajor)
	assert.Equal(t, expectedHTTPRequest.ProtoMinor, actuallyHTTPRequest.ProtoMinor)
	assert.Equal(t, expectedHTTPRequest.Header, actuallyHTTPRequest.Header)
	assert.Equal(t, expectedHTTPRequest.ContentLength, actuallyHTTPRequest.ContentLength)
	assert.Equal(t, expectedHTTPRequest.TransferEncoding, actuallyHTTPRequest.TransferEncoding)
	assert.Equal(t, expectedHTTPRequest.Close, actuallyHTTPRequest.Close)
	assert.Equal(t, expectedHTTPRequest.Form, actuallyHTTPRequest.Form)
	assert.Equal(t, expectedHTTPRequest.PostForm, actuallyHTTPRequest.PostForm)
	assert.Equal(t, expectedHTTPRequest.MultipartForm, actuallyHTTPRequest.MultipartForm)
	assert.Equal(t, expectedHTTPRequest.Trailer, actuallyHTTPRequest.Trailer)
	assert.Equal(t, expectedHTTPRequest.RemoteAddr, actuallyHTTPRequest.RemoteAddr)
	assert.Equal(t, expectedHTTPRequest.RequestURI, actuallyHTTPRequest.RequestURI)
	assert.Equal(t, expectedHTTPRequest.TLS, actuallyHTTPRequest.TLS)
	assert.Equal(t, expectedHTTPRequest.Cancel, actuallyHTTPRequest.Cancel)
	assert.Equal(t, expectedHTTPRequest.Response, actuallyHTTPRequest.Response)

	eBodyBuffer, _ := expectedHTTPRequest.GetBody()
	expectedBody, _ := ioutil.ReadAll(eBodyBuffer)

	aBodyBuffer, _ := actuallyHTTPRequest.GetBody()
	actuallyBody, _ := ioutil.ReadAll(aBodyBuffer)

	assert.Equal(t, expectedBody, actuallyBody)

}
