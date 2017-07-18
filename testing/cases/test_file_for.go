package cases

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func testFilesFor(t testing.TB, dirname string, includes, ignores []string) []string {
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

		if !isIncludedTestCase(file.Name(), includes, ignores) {
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

func isIncludedTestCase(s string, includes, ignores []string) bool {

	for _, ignore := range ignores {
		rgx := rgxForMatcher(ignore)

		if rgx.MatchString(s) {
			return false
		}
	}

	for _, include := range includes {
		rgx := rgxForMatcher(include)

		if !rgx.MatchString(s) {
			return false
		}
	}

	return true

}

func rgxForMatcher(s string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + regexp.QuoteMeta(s))
}
