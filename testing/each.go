package testing

import (
	"testing"

	escher "github.com/adamluzsi/escher-go"
)

func EachTestConfigFor(t testing.TB, topic string, tester func(escher.Config, TestConfig) bool) {
	testedCases := make(map[bool]struct{})

	for _, testConfig := range getTestConfigsForTopic(t, topic) {
		t.Log("-----------------------------------------------")
		t.Log(testConfig.getTitle())
		if testConfig.Description != "" {
			t.Log(testConfig.Description)
		}

		testedCases[tester(fixedConfigBy(t, testConfig.Config), testConfig)] = struct{}{}
		t.Log("-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --")
	}

	if _, ok := testedCases[true]; !ok {
		t.Fatal("No test case was used")
	}
}
