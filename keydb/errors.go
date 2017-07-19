package keydb

import "errors"

var (
	KeyIDNotFound          = errors.New("KeyID Not Found")
	KeyPoolEnvValueIsEmpty = errors.New("KEY_POOL Env value is empty")
)
