package testing

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func testFilesFor(t testing.TB, dirname, topic string) []string {
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

		rgx := regexp.MustCompile(regexp.QuoteMeta(formattedPrefix(topic)))

		if !rgx.MatchString(file.Name()) {
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

func formattedPrefix(prefix string) string {
	return strings.ToLower(prefix)
}
