package debug

import "os"

var enabled bool

func init() {
	_, isSet := os.LookupEnv("ESCHER_DEBUG")

	enabled = isSet
}
