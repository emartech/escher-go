package cases

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/debug"
)

func EachTestConfigFor(t *testing.T, includes, ignores []string, tester func(*testing.T, config.Config, TestConfig) bool) {
	testedCases := make(map[bool]struct{})

	for _, testConfig := range getTestConfigsForTopic(t, includes, ignores) {

		testCaseName := path.Join(
			// this works nicely because the test-cases project
			// has directory level for each test case group
			filepath.Base(filepath.Dir(testConfig.FilePath)),
			filepath.Base(testConfig.FilePath),
		)

		t.Run(testCaseName, func(t *testing.T) {
			t.Log(testConfig.getTitle())
			t.Log(testConfig.Description)
			t.Log(testConfig.FilePath)
			defer debug.SetPrinter(t.Log)()

			testedCases[tester(t, fixedConfigBy(testConfig.Config), testConfig)] = struct{}{}
		})
	}

	if _, ok := testedCases[true]; !ok {
		t.Fatal("No test case matched the current include/ignore setup")
	}
}
