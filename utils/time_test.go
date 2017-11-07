package utils_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/EscherAuth/escher/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseTimeValidTimeStringGiven(t *testing.T) {

	// referenceTime, _ := time.Parse(time.RFC3339Nano, time.RFC3339Nano)
	referenceTime := time.Now()

	fmt.Println(time.Now().Format("20060102T150405Z07"))

	return
	supportedFormats := []string{
		"20060102T150405Z07",
		"20060102T150405Z0700",
		"20060102T150405Z07:00",
		"20060102T150405Z070000",
		"20060102T150405Z07:00:00",
		"Fri, 02 Jan 2006 15:04:05 MST",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	for _, format := range supportedFormats {
		timeString := referenceTime.Format(format)

		expectedTime, _ := time.Parse(format, timeString)
		actuallyTime, err := utils.ParseTime(timeString)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedTime, actuallyTime)
	}

}
