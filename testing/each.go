package testing

import (
	"testing"

	escher "github.com/EscherAuth/escher"
)

func EachTestConfigFor(t testing.TB, topic string, tester func(escher.Config, TestConfig) bool) {
	testedCases := make(map[bool]struct{})

	for _, testConfig := range getTestConfigsForTopic(t, topic) {
		testedCases[tester(fixedConfigBy(t, testConfig.Config), testConfig)] = struct{}{}

		if testing.Verbose() {
			t.Log("-----------------------------------------------")

			t.Log(testConfig.getTitle())
			t.Log(testConfig.FilePath)

			if testConfig.Description != "" {
				t.Log(testConfig.Description)
			}

			if t.Failed() {
				if isFastFailEnabled() {
					t.FailNow()
				}
			} else {
				t.Log("OK")
			}

			t.Log("-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --\n")
		}

	}

	if _, ok := testedCases[true]; !ok {
		t.Fatal("No test case was used")
	}
}
