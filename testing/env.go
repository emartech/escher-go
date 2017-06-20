package testing

import (
	"os"
	"strings"
)

func isFastFailEnabled() bool {
	ff := os.Getenv("FAST_FAIL")

	return strings.ToLower(ff) == "true"
}
