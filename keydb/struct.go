package keydb

import "errors"

type KeyDB interface {
	GetSecret(keyID string) (string, error)
}

type keydb struct {
	db map[string]string
}

func NewBySlice(raw [][2]string) KeyDB {
	kd := &keydb{make(map[string]string)}

	for _, keyPairs := range raw {
		kd.db[keyPairs[0]] = keyPairs[1]
	}

	return kd
}

var KeyIDNotFound = errors.New("KeyID Not Found")

func (kd *keydb) GetSecret(keyID string) (string, error) {
	secret, ok := kd.db[keyID]

	if !ok {
		return secret, KeyIDNotFound
	}

	return secret, nil
}
