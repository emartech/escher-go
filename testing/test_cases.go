package testing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	escher "github.com/adamluzsi/escher-go"
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
		// "presignUrl": []string{
		// 	"presignurl-valid-with-path-query",
		// },
		"authenticate": []string{
			"authenticate-valid-authentication-datein-expiretime",
			"authenticate-valid-get-vanilla-empty-query",
			"authenticate-valid-get-vanilla-empty-query-with-custom-headernames",
			"authenticate-valid-presigned-url-with-query",
			"authenticate-valid-ignore-headers-order",
			"authenticate-error-host-header-not-signed",
			"authenticate-error-date-header-not-signed",
			"authenticate-error-invalid-auth-header",
			"authenticate-error-invalid-escher-key",
			"authenticate-error-invalid-credential-scope",
			"authenticate-error-invalid-hash-algorithm",
			"authenticate-error-missing-auth-header",
			"authenticate-error-missing-host-header",
			"authenticate-error-missing-date-header",
			"authenticate-error-date-header-auth-header-date-not-equal",
			"authenticate-error-request-date-invalid",
			"authenticate-error-wrong-signature",
			"authenticate-error-presigned-url-expired",
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
	Request              escher.Request `json:"request"`
	CanonicalizedRequest string         `json:"canonicalizedRequest"`
	StringToSign         string         `json:"stringToSign"`
	AuthHeader           string         `json:"authHeader"`
	Error                string         `json:"error"`
}

type TestConfig struct {
	ID            string
	HeadersToSign []string           `json:"headersToSign"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Request       escher.Request     `json:"request"`
	Config        escher.Config      `json:"config"`
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
	// TODO: fix this rel path

	var filename = filepath.Join(testSuitePath(t), testSuite+"_testsuite", testID+".json")
	fmt.Println(filename)
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		t.Fatal(err)
	}

	json.Unmarshal(content, &testConfig)
	testConfig.ID = testSuite + ":" + testID
	return testConfig
}

func testSuitePath(t testing.TB) string {
	testSuitePath := os.Getenv("TEST_SUITE_PATH")

	if testSuitePath == "" {
		t.Fatal("TEST_SUITE_PATH env is missing, can't find the escher tests")
	}

	_, err := os.Stat(testSuitePath)

	if err != nil && os.IsNotExist(err) {
		t.Fatal("given TEST_SUITE_PATH IsNotExists!")
	}

	return testSuitePath
}

func EachTestConfigFor(t testing.TB, topic string, tester func(escher.Config, TestConfig) bool) {
	testedCases := make(map[bool]struct{})

	for _, testConfig := range getTestConfigsForTopic(t, topic) {
		config := fixedConfigBy(testConfig.Config)
		t.Log(testConfig.getTitle())
		t.Log(testConfig.Description)
		testedCases[tester(config, testConfig)] = struct{}{}
	}

	if _, ok := testedCases[true]; !ok {
		t.Fatal("No test case was used")
	}
}

func fixedConfigBy(config escher.Config) escher.Config {
	var t, err = time.Parse("2006-01-02T15:04:05.999999Z", config.Date)
	if err != nil {
		t, err = time.Parse("Fri, 02 Jan 2006 15:04:05 GMT", config.Date)
	}
	if err != nil {
		t = time.Now().UTC()
	}
	config.Date = t.Format("20060102T150405Z")
	return config
}
