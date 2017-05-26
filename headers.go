package escher

type RequestHeaders [][2]string

func (rh RequestHeaders) Get(key string) (string, bool) {
	for _, keyPairs := range rh {
		if keyPairs[0] == key {
			return keyPairs[1], true
		}
	}
	return "", false
}