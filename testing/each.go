package testing

import (
	"testing"

	escher "github.com/EscherAuth/escher"
)

func EachTestConfigFor(t testing.TB, topic string, tester func(escher.Config, TestConfig) bool) {
	testedCases := make(map[bool]struct{})

	for _, testConfig := range getTestConfigsForTopic(t, topic) {

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
			}

			t.Log("-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --\n")
		}

		if isFastFailEnabled() {
			t.FailNow()
		}

	}

	if _, ok := testedCases[true]; !ok {
		t.Fatal("No test case was used")
	}
}
