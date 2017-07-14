package keydb

func (kd *keydb) GetSecret(keyID string) (string, error) {
	secret, ok := kd.db[keyID]

	if !ok {
		return secret, KeyIDNotFound
	}

	return secret, nil
}
