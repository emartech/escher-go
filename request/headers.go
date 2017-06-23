package request

type Headers [][2]string

func (h Headers) Get(key string) (string, bool) {
	for _, keyPairs := range h {
		if keyPairs[0] == key {
			return keyPairs[1], true
		}
	}
	return "", false
}
