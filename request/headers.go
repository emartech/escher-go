package request

import (
	"strings"
)

type Headers [][2]string

func (h Headers) Get(key string) (string, bool) {
	expectedKeyToFind := strings.ToLower(key)

	for _, keyPairs := range h {
		if strings.ToLower(keyPairs[0]) == expectedKeyToFind {
			return keyPairs[1], true
		}
	}

	return "", false
}
