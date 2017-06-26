package testing

import (
	"flag"
	"os"
	"regexp"
	"strings"
)

var fastFail bool

func init() {
	const description = "set the test to exit on the first fail in every test case"
	flag.BoolVar(&fastFail, "fast-fail", false, description)
	flag.BoolVar(&fastFail, "ff", false, description+" (shorthand)")
}

func isFastFailEnabled() bool {
	return strings.ToLower(os.Getenv("FAST_FAIL")) == "true" || fastFail
}

var testCase = flag.String("case", "", "run only the following test case")

func isTestCaseAllowed(testCaseFileName string) bool {
	if *testCase == "" {
		return true
	}

	rgx := regexp.MustCompile(regexp.QuoteMeta(*testCase))
	return rgx.MatchString(testCaseFileName)
}
