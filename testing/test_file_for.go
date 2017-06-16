package testing

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func testFilesFor(t testing.TB, dirname, prefix string) []string {
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

		if !strings.HasPrefix(file.Name(), formattedPrefix(prefix)) {
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
