package validator

import (
	escher "github.com/adamluzsi/escher-go"
	"github.com/adamluzsi/escher-go/keydb"
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
