package request_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/signer"
	"github.com/EscherAuth/escher/testing/cases"
)

func TestUpdateHTTPRequest(t *testing.T) {
	t.Parallel()

	cases.EachTestConfigFor(t, []string{"signrequest", "post", "payload", "utf8"}, nil, func(c config.Config, tc cases.TestConfig) bool {

		HTTPRequestToUpdate, err := tc.Request.HTTPRequest("http://example.org:8080")

		if err != nil {
			t.Fatal(err)
		}

		SignedRequest, err := signer.New(c).SignRequest(&tc.Request, tc.HeadersToSign)

		if err != nil {
			t.Fatal(err)
		}

		mergeErr := SignedRequest.UpdateHTTPRequest(HTTPRequestToUpdate)

		if mergeErr != nil {
			t.Fatal(mergeErr)
		}

		assert.Equal(t, "http://example.org:8080/", HTTPRequestToUpdate.URL.String())

		for _, KeyValues := range SignedRequest.Headers() {
			assert.Equal(t, KeyValues[1], HTTPRequestToUpdate.Header.Get(KeyValues[0]))
		}

		return !t.Failed()
	})

}
