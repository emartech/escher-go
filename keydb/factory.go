package keydb

import (
	"os"
)

func NewByKeyValuePair(apiKey, secret string) KeyDB {
	return NewBySlice([][2]string{[2]string{apiKey, secret}})
}

func NewBySlice(raw [][2]string) KeyDB {
	kd := make(keydb)

	for _, keyPairs := range raw {
		kd[keyPairs[0]] = append(kd[keyPairs[0]], keyPairs[1])
	}

	return kd
}

func NewFromENV() (KeyDB, error) {
	jsonString, isGiven := os.LookupEnv("KEY_POOL")

	if !isGiven {
		return nil, KeyPoolEnvValueIsEmpty
	}

	return parseFromKeyPool(jsonString)
}
