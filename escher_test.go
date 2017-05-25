package escher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = map[string]map[string][]string{
	"aws4": map[string][]string{
		"signRequest": []string{
			"signrequest-get-vanilla",
			"signrequest-post-vanilla",
			"signrequest-get-vanilla-query",
			"signrequest-post-vanilla-query",
			"signrequest-get-vanilla-empty-query-key",
			"signrequest-post-vanilla-empty-query-value",
			"signrequest-get-vanilla-query-order-key",
			"signrequest-post-x-www-form-urlencoded",
			"signrequest-post-x-www-form-urlencoded-parameters",
			"signrequest-get-header-value-trim",
			"signrequest-post-header-key-case",
			"signrequest-post-header-key-sort",
			"signrequest-post-header-value-case",
			"signrequest-get-vanilla-query-order-value",
			"signrequest-get-vanilla-query-order-key-case",
			"signrequest-get-unreserved",
			"signrequest-get-vanilla-query-unreserved",
			"signrequest-get-vanilla-ut8-query",
			"signrequest-get-utf8",
			"signrequest-get-space",
			"signrequest-post-vanilla-query-space",
			"signrequest-post-vanilla-query-nonunreserved",
			"signrequest-get-slash",
			"signrequest-get-slashes",
			"signrequest-get-slash-dot-slash",
			"signrequest-get-slash-pointless-dot",
			"signrequest-get-relative",
			"signrequest-get-relative-relative",
		},
		"presignUrl":   []string{},
		"authenticate": []string{},
	},
	"emarsys": map[string][]string{
		"signRequest": []string{
			"signrequest-get-header-key-duplicate",
			"signrequest-get-header-value-order",
			"signrequest-post-header-key-order",
			"signrequest-post-header-value-spaces",
			"signrequest-post-header-value-spaces-within-quotes",
			"signrequest-post-payload-utf8",
			"signrequest-date-header-should-be-signed-headers",
			"signrequest-support-custom-config",
			"signrequest-only-sign-specified-headers",
		},
		"presignUrl": []string{
			"presignurl-valid-with-path-query",
		},
		"authenticate": []string{
		// "authenticate-valid-authentication-datein-expiretime",
		// "authenticate-valid-get-vanilla-empty-query",
		// "authenticate-valid-get-vanilla-empty-query-with-custom-headernames",
		// "authenticate-valid-presigned-url-with-query",
		// "authenticate-valid-ignore-headers-order",
		// "authenticate-error-host-header-not-signed",
		// "authenticate-error-date-header-not-signed",
		// "authenticate-error-invalid-auth-header",
		// "authenticate-error-invalid-escher-key",
		// "authenticate-error-invalid-credential-scope",
		// "authenticate-error-invalid-hash-algorithm",
		// "authenticate-error-missing-auth-header",
		// "authenticate-error-missing-host-header",
		// "authenticate-error-missing-date-header",
		// "authenticate-error-date-header-auth-header-date-not-equal",
		// "authenticate-error-request-date-invalid",
		// "authenticate-error-wrong-signature",
		// "authenticate-error-presigned-url-expired",
		},
	},
}

func getTestConfigsForTopic(t testing.TB, topic string) []TestConfig {
	var configs = []TestConfig{}
	for testSuite, testTypes := range tests {
		for testTopic, testFiles := range testTypes {
			if testTopic == topic {
				for _, testFile := range testFiles {
					configs = append(configs, loadTestFile(t, testSuite, testFile))
				}
			}
		}
	}
	return configs
}

type TestConfigExpected struct {
	Request              EscherRequest `json:"request"`
	CanonicalizedRequest string        `json:"canonicalizedRequest"`
	StringToSign         string        `json:"stringToSign"`
	AuthHeader           string        `json:"authHeader"`
}

type TestConfig struct {
	ID            string
	HeadersToSign []string           `json:"headersToSign"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Request       EscherRequest      `json:"request"`
	Config        EscherConfig       `json:"config"`
	Expected      TestConfigExpected `json:"expected"`
}

func (testConfig TestConfig) getTitle() string {
	var title string
	if testConfig.Title == "" {
		title = testConfig.ID
	} else {
		title = testConfig.Title + " (" + testConfig.ID + ")"
	}
	return title
}

func loadTestFile(t testing.TB, testSuite string, testID string) TestConfig {
	if testing.Verbose() {
		log.Printf("%s - %s\n", testSuite, testID)
	}

	var testConfig TestConfig
	var filename = testSuite + "_testsuite/" + testID + ".json"
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		t.Fatal(err)
	}

	json.Unmarshal(content, &testConfig)
	testConfig.ID = testSuite + ":" + testID
	return testConfig
}

func eachTestConfigFor(t testing.TB, topic string, tester func(EscherConfig, TestConfig)) {
	for _, testConfig := range getTestConfigsForTopic(t, topic) {
		var escher = Escher(testConfig.Config)
		t.Log(testConfig.getTitle())
		t.Log(testConfig.Description)
		tester(escher, testConfig)
	}
}

func TestCanonicalizeRequest(t *testing.T) {
	var tested bool
	t.Log("CanonicalizeRequest should return with a proper string")
	eachTestConfigFor(t, "signRequest", func(escher EscherConfig, testConfig TestConfig) {
		if testConfig.Expected.CanonicalizedRequest == "" {
			return
		}

		tested = true

		canonicalizedRequest := escher.CanonicalizeRequest(testConfig.Request, testConfig.HeadersToSign)

		if canonicalizedRequest != testConfig.Expected.CanonicalizedRequest {
			t.Fatal("canonicalizedRequest should be eq")
		}

	})
	if !tested {
		t.Fatal("TestCanonicalizeRequest not tested")
	}
}

func TestGetStringToSign(t *testing.T) {
	var tested bool
	t.Log("GetStringToSign should return with a proper string")
	eachTestConfigFor(t, "signRequest", func(escher EscherConfig, testConfig TestConfig) {
		if testConfig.Expected.StringToSign == "" {
			return
		}

		tested = true

		stringToSign := escher.GetStringToSign(testConfig.Request, testConfig.HeadersToSign)
		if stringToSign != testConfig.Expected.StringToSign {
			t.Fatal("stringToSign expected to eq with the test config expectation")
		}
	})
	if !tested {
		t.Fatal("TestGetStringToSign not tested")
	}
}

func TestGenerateHeader(t *testing.T) {
	var tested bool

	t.Log("GenerateHeader should return with a proper string")
	eachTestConfigFor(t, "signRequest", func(escher EscherConfig, testConfig TestConfig) {
		if testConfig.Expected.AuthHeader == "" {
			return
		}

		tested = true
		authHeader := escher.GenerateHeader(testConfig.Request, testConfig.HeadersToSign)
		if authHeader != testConfig.Expected.AuthHeader {
			t.Fatal("authHeader generation failed")
		}
	})

	if !tested {
		t.Fatal("TestGenerateHeader not tested")
	}
}

func TestSignRequest(t *testing.T) {
	// t.Skip("no use case that would test this feature")
	var tested bool
	t.Log("SignRequest should return with a properly signed request")

	eachTestConfigFor(t, "signRequest", func(escher EscherConfig, testConfig TestConfig) {
		if testConfig.Expected.Request.Method == "" {
			return
		}

		tested = true
		request := escher.SignRequest(testConfig.Request, testConfig.HeadersToSign)
		assert.Equal(t, testConfig.Expected.Request, request, "Requests should be eq")
	})

	if !tested {
		t.Fatal("TestSignRequest not tested")
	}
}

// func TestAuthenticate(t *testing.T) {
// 	for _, testConfig := range getTestConfigsForTopic(t, "authenticate") {
// 		var escher = Escher(testConfig.Config)
// 		var testTitle = testConfig.getTitle()

// 		if testConfig.Expected.CanonicalizedRequest != "" {

// 			var canonicalizedRequest = escher.CanonicalizeRequest(testConfig.Request, testConfig.HeadersToSign)
// 			So(canonicalizedRequest, ShouldEqual, testConfig.Expected.CanonicalizedRequest)
// 		}
// 	}
// }
