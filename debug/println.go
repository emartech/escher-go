package debug

import (
	"fmt"
)

// Println for debugging purpose
func Println(a ...interface{}) (n int, err error) {

	if !enabled {
		return 0, nil
	}

	return fmt.Println(a...)

}
