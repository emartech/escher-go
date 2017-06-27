package testing

import (
	"testing"

	escher "github.com/EscherAuth/escher"
)

func EachTestConfigFor(t testing.TB, topics []string, tester func(escher.Config, TestConfig) bool) {
	testedCases := make(map[bool]struct{})

	for _, testConfig := range getTestConfigsForTopic(t, topics) {

		if testing.Verbose() {
			t.Log("-----------------------------------------------")

			t.Log(testConfig.getTitle())
			if testConfig.Description != "" {
				t.Log(testConfig.Description)
			}

			t.Log(testConfig.FilePath)
		}

		testedCases[tester(fixedConfigBy(t, testConfig.Config), testConfig)] = struct{}{}

		if testing.Verbose() {

			if !t.Failed() {
				t.Log("OK")
			} else {
				t.Log("ERROR")
			}

			t.Log("-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --\n")
		}

		if isFastFailEnabled() && t.Failed() {
			t.FailNow()
		}

	}

	if _, ok := testedCases[true]; !ok {
		t.Fatal("No test case was used")
	}
}
