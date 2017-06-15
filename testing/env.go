package testing

import (
	"os"
	"strings"
)

func isFailFastEnabled() bool {
	ff := os.Getenv("FAIL_FAST")

	return strings.ToLower(ff) == "true"
}
