package mock

import (
	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
)

type ReceivedValidations []ReceivedValidation

type ReceivedValidation struct {
	Request          request.Interface
	KeyDB            keydb.KeyDB
	MandatoryHeaders []string
}
