package validator

import (
	escher "github.com/EscherAuth/escher"
	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
)

type Validator interface {
	Validate(request.Interface, keydb.KeyDB, []string) (string, error)
}

type validator struct {
	config escher.Config
}

func New(config escher.Config) Validator {
	return &validator{config}
}
