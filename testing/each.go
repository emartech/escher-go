package testing

import (
	"testing"

	escher "github.com/EscherAuth/escher"
)

func EachTestConfigFor(t testing.TB, topic string, tester func(escher.Config, TestConfig) bool) {
	testedCases := make(map[bool]struct{})

	for _, testConfig := range getTestConfigsForTopic(t, topic) {

		testedCases[tester(fixedConfigBy(t, testConfig.Config), testConfig)] = struct{}{}

		if t.Failed() || testing.Verbose() {
			t.Log("-----------------------------------------------")
			t.Log(testConfig.getTitle())
			t.Log(testConfig.FilePath)
			if testConfig.Description != "" {
				t.Log(testConfig.Description)
			}
			t.Log("-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --")

			if isFailFastEnabled() {
				t.FailNow()
			}
		}

	}

	if _, ok := testedCases[true]; !ok {
		t.Fatal("No test case was used")
	}
}
