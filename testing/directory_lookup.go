package testing

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

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

		if file.Name() == "ducktype_cases" {
			continue
		}

		dirs = append(dirs, filepath.Join(testDir, file.Name()))

	}

	return dirs
}
