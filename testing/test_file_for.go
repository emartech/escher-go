package testing

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func testFilesFor(t testing.TB, dirname string, topics []string) []string {
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

		if !MatchAll(file.Name(), topics) {
			continue
		}

		if !isTestCaseAllowed(file.Name()) {
			continue
		}

		if !file.IsDir() {
			testFiles = append(testFiles, filepath.Join(dirname, file.Name()))
		}

	}

	return testFiles
}

func MatchAll(s string, matchers []string) bool {
	for _, matcher := range matchers {
		rgx := regexp.MustCompile(regexp.QuoteMeta(strings.ToLower(matcher)))

		if !rgx.MatchString(s) {
			return false
		}
	}
	return true
}
