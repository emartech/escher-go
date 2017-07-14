package mock_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
	"github.com/EscherAuth/escher/validator"
)

func requestBy(t testing.TB, method, path, body string) *request.Request {

	requestBody := bytes.NewReader([]byte(body))
	httpRequest, err := http.NewRequest(method, path, requestBody)

	if err != nil {
		t.Fatal(err)
	}

	return request.NewFromHTTPRequest(httpRequest)
}

func useValidator(t testing.TB, v validator.Validator, request request.Interface, keydb keydb.KeyDB, mandatoryHeaders []string) (string, error) {
	return v.Validate(request, keydb, mandatoryHeaders)
}
