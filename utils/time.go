package utils

import (
	"errors"
	"time"
)

const EscherDateFormat = "20060102T150405Z"

var acceptedTimeFormats = []string{
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
	EscherDateFormat,
}

func ParseTime(timeStr string) (time.Time, error) {
	for _, layout := range acceptedTimeFormats {
		t, err := time.Parse(layout, timeStr)

		if err == nil {
			return t, err
		}
	}

	return time.Time{}, errors.New("no layout found for " + timeStr)
}
