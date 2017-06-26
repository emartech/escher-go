package testing

import (
	"flag"
	"os"
	"regexp"
	"strings"
)

var fastFail = flag.Bool("fast-fail", false, "set the test to exit on the first fail in every test case")

func isFastFailEnabled() bool {
	return strings.ToLower(os.Getenv("FAST_FAIL")) == "true" || *fastFail
}

var testCase = flag.String("case", "", "run only the following test case")

func isTestCaseAllowed(testCaseFileName string) bool {
	if *testCase == "" {
		return true
	}

	rgx := regexp.MustCompile(*testCase)
	return rgx.MatchString(testCaseFileName)
}
