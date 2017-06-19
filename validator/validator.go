package validator

import (
	escher "github.com/EscherAuth/escher"
	"github.com/EscherAuth/escher/keydb"
)

type Validator interface {
	Validate(request escher.Request, keyDB keydb.KeyDB, mandatoryHeaders []string) (string, error)
}

type validator struct {
	config escher.Config
}

func New(config escher.Config) Validator {
	return &validator{config}
}
