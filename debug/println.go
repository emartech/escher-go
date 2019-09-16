package debug

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var printer func(...interface{})

func SetPrinter(newPrinter func(...interface{})) (teardown func()) {
	currentPrinter := printer
	printer = newPrinter
	return func() { printer = currentPrinter }
}

const warningMessage = `                                                                              
WARNING: ESCHER_DEBUG is enabled!                                                                     
This should be only used in development for escher related problems debugging.                        
With this sensitive information could be logged to the STDOUT!                                        
Use with extra caution                                                                                
`

func init() {
	printer = func(...interface{}) {}

	v, isSet := os.LookupEnv("ESCHER_DEBUG")
	if isSet && strings.ToUpper(v) == "TRUE" {
		_, _ = fmt.Fprint(os.Stderr, warningMessage)
		printer = log.New(os.Stdout, `escher`, log.LstdFlags).Println
	}
}

// Println for debugging purpose
func Println(v ...interface{}) {
	printer(v...)
}

