package keydb

func (kd keydb) GetSecret(keyID string) (string, error) {
	secrets, ok := kd[keyID]

	if !ok {
		return "", KeyIDNotFound
	}

	return secrets[0], nil
}
