package keydb

func NewByKeyValuePair(apiKey, secret string) KeyDB {
	return NewBySlice([][2]string{[2]string{apiKey, secret}})
}

func NewBySlice(raw [][2]string) KeyDB {
	kd := &keydb{make(map[string]string)}

	for _, keyPairs := range raw {
		kd.db[keyPairs[0]] = keyPairs[1]
	}

	return kd
}
