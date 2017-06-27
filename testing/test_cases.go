package testing

import (
	"os"
	"testing"
)

func getTestConfigsForTopic(t testing.TB, includes, ignores []string) []TestConfig {
	configs := make([]TestConfig, 0)

	for _, dirPath := range getTestCaseDirectories(t) {
		for _, filePath := range testFilesFor(t, dirPath, includes, ignores) {
			configs = append(configs, testConfigBy(t, filePath))
		}
	}

	return configs
}

func testSuitePath(t testing.TB) string {
	testSuitePath := os.Getenv("TEST_CASES_PATH")

	if testSuitePath == "" {
		t.Fatal("TEST_CASES_PATH env is missing, can't find the escher tests")
	}

	_, err := os.Stat(testSuitePath)

	if err != nil && os.IsNotExist(err) {
		t.Fatal("given TEST_CASES_PATH IsNotExists!")
	}

	return testSuitePath
}
