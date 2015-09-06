package escher

import (
  "testing"
  . "github.com/smartystreets/goconvey/convey"
  "io/ioutil"
  "encoding/json"
)

var tests = map[string]map[string][]string {
  "aws4": map[string][]string {
    "signRequest": []string {
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
    "presignUrl": []string {},
    "authenticate": []string {},
  },
  "emarsys": map[string][]string {
    "signRequest": []string {
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
    "presignUrl": []string {
      // "presignurl-valid-with-path-query",
    },
    "authenticate": []string {
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
};

func getTestConfigsForTopic(topic string) []TestConfig {
  var configs = []TestConfig {};
  for testSuite, testTypes := range tests {
    for testTopic, testFiles := range testTypes {
      if testTopic == topic {
        for _, testFile := range testFiles {
          configs = append(configs, loadTestFile(testSuite, testFile));
        }
      }
    }
  }
  return configs
};

type TestConfigExpected struct {
  Request EscherRequest
  CanonicalizedRequest string
  StringToSign string
  AuthHeader string
}

type TestConfig struct {
  ID string
  HeadersToSign []string
  Title string
  Description string
  Request EscherRequest
  Config EscherConfig
  Expected TestConfigExpected
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

func loadTestFile(testSuite string, testID string) TestConfig {
  var testConfig TestConfig
  var filename = testSuite + "_testsuite/" + testID + ".json"
  content, _ := ioutil.ReadFile(filename)
  json.Unmarshal(content, &testConfig)
  testConfig.ID = testSuite + ":" + testID
  return testConfig
}

func TestCanonicalizeRequest(t *testing.T) {

  Convey("CanonicalizeRequest should return with a proper string", t, func() {
    for _, testConfig := range getTestConfigsForTopic("signRequest") {
      var escher = Escher(testConfig.Config)
      var testTitle = testConfig.getTitle()
      if testConfig.Expected.CanonicalizedRequest != "" {
        Convey(testTitle, func() {
          var canonicalizedRequest = escher.CanonicalizeRequest(testConfig.Request, testConfig.HeadersToSign)
          So(canonicalizedRequest, ShouldEqual, testConfig.Expected.CanonicalizedRequest)
        })
      }
    }
  })

  Convey("GetStringToSign should return with a proper string", t, func() {
    for _, testConfig := range getTestConfigsForTopic("signRequest") {
      var escher = Escher(testConfig.Config)
      var testTitle = testConfig.getTitle()
      if testConfig.Expected.StringToSign != "" {
        Convey(testTitle, func() {
          var stringToSign = escher.GetStringToSign(testConfig.Request, testConfig.HeadersToSign)
          So(stringToSign, ShouldEqual, testConfig.Expected.StringToSign)
        })
      }
    }
  })

  Convey("GenerateHeader should return with a proper string", t, func() {
    for _, testConfig := range getTestConfigsForTopic("signRequest") {
      var escher = Escher(testConfig.Config)
      var testTitle = testConfig.getTitle()
      if testConfig.Expected.AuthHeader != "" {
        Convey(testTitle, func() {
          var authHeader = escher.GenerateHeader(testConfig.Request, testConfig.HeadersToSign)
          So(authHeader, ShouldEqual, testConfig.Expected.AuthHeader)
        })
      }
    }
  })

  Convey("SignRequest should return with a properly signed request", t, func() {
    for _, testConfig := range getTestConfigsForTopic("signRequest") {
      var escher = Escher(testConfig.Config)
      var testTitle = testConfig.getTitle()
      if testConfig.Expected.Request.Method != "" {
        Convey(testTitle, func() {
          var request = escher.SignRequest(testConfig.Request, testConfig.HeadersToSign)
          So(request, ShouldResemble, testConfig.Expected.Request)
        })
      }
    }
  })

}
