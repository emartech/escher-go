package utils

import (
	"errors"
	"time"
)

const EscherDateFormat = "20060102T150405Z07"
const HTTPHeaderFormat = "Mon, 02 Jan 2006 15:04:05 MST"

var acceptedTimeFormats = []string{
	EscherDateFormat,
	"20060102T150405Z0700",
	"20060102T150405Z07:00",
	"20060102T150405Z070000",
	"20060102T150405Z07:00:00",
	HTTPHeaderFormat,
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

func ParseTime(timeStr string) (time.Time, error) {
	for _, layout := range acceptedTimeFormats {
		t, err := time.Parse(layout, timeStr)

		if err == nil {
			return t, err
		}
	}

	return time.Time{}, errors.New("no layout found for " + timeStr)
}
