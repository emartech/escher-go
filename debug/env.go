package debug

import (
	"fmt"
	"os"
)

var enabled bool

const warningMessage = `
WARNING: ESCHER_DEBUG is enabled!
This should be only used in development for escher related problems debugging.
With this sensitive information could be logged to the STDOUT!
Use with extra caution
`

func init() {
	_, isSet := os.LookupEnv("ESCHER_DEBUG")

	if isSet {
		fmt.Println(warningMessage)
	}

	enabled = isSet
}
