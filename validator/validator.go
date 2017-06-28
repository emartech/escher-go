package validator

import (
	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
)

type Validator interface {
	Validate(request.Interface, keydb.KeyDB, []string) (string, error)
}

type validator struct {
	config config.Config
}

func New(c config.Config) Validator {
	return &validator{c}
}
