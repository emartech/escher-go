package keydb

type KeyDB interface {
	GetSecret(keyID string) (string, error)
}

type keydb struct {
	db map[string]string
}
