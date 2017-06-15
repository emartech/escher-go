package testing

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func getTestConfigsForTopic(t testing.TB, topic string) []TestConfig {
	configs := make([]TestConfig, 0)

	for _, dirPath := range getTestCaseDirectories(t) {
		for _, filePath := range testFilesFor(t, dirPath, topic) {
			configs = append(configs, testConfigBy(t, filePath))
		}
	}

	return configs
}

func testSuitePath(t testing.TB) string {
	testSuitePath := os.Getenv("TEST_SUITE_PATH")

	if testSuitePath == "" {
		t.Fatal("TEST_SUITE_PATH env is missing, can't find the escher tests")
	}

	_, err := os.Stat(testSuitePath)

	if err != nil && os.IsNotExist(err) {
		t.Fatal("given TEST_SUITE_PATH IsNotExists!")
	}

	return testSuitePath
}

func getTestCaseDirectories(t testing.TB) []string {
	testDir := testSuitePath(t)

	files, err := ioutil.ReadDir(testDir)

	if err != nil {
		t.Fatal(err)
	}

	dirs := make([]string, 0)

	for _, file := range files {

		if !file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		dirs = append(dirs, filepath.Join(testDir, file.Name()))

	}

	return dirs
}

func testFilesFor(t testing.TB, dirname, suffix string) []string {
	files, err := ioutil.ReadDir(dirname)

	if err != nil {
		t.Fatal(err)
	}

	testFiles := make([]string, 0)

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		if !strings.HasSuffix(file.Name(), suffix) {

		}

		if !file.IsDir() {
			testFiles = append(testFiles, filepath.Join(dirname, file.Name()))
		}

	}

	return testFiles
}
