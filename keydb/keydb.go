package keydb

type KeyDB interface {
	GetSecret(keyID string) (string, error)
}

type keydb map[string][]string
