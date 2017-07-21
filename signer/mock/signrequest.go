package mock

import (
	"github.com/EscherAuth/escher/request"
)

type signRequestResponse struct {
	Request *request.Request
	Error   error
}

func (m *Mock) AddSignRequestResponse(req *request.Request, err error) {
	newSignRequestResponse := signRequestResponse{Request: req, Error: err}
	m.signRequestResponses = append(m.signRequestResponses, newSignRequestResponse)
}

func (m *Mock) nextSignRequestResponse() (*request.Request, error) {
	first := m.signRequestResponses[0]
	m.signRequestResponses = m.signRequestResponses[1:]
	return first.Request, first.Error
}

func (m *Mock) SignRequest(r request.Interface, headersToSign []string) (*request.Request, error) {
	if len(m.signRequestResponses) == 0 {
		return r.(*request.Request), nil
	} else {
		return m.nextSignRequestResponse()
	}
}
