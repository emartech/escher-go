package mock

import (
	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
)

func (m *Mock) Validate(escherRequest request.Interface, db keydb.KeyDB, mandatoryHeaders []string) (string, error) {

	toBeValidated := ReceivedValidation{
		Request:          escherRequest,
		KeyDB:            db,
		MandatoryHeaders: mandatoryHeaders,
	}

	m.ReceivedValidations = append(m.ReceivedValidations, toBeValidated)

	return m.popValidationResults()

}
