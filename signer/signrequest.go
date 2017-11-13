package signer

import "github.com/EscherAuth/escher/request"

func (s *signer) SignRequest(r request.Interface, headersToSign []string) (*request.Request, error) {
	err := s.config.Validate(r)

	if err != nil {
		return nil, err
	}

	headers := r.Headers()

	var authHeader = s.generateHeader(r, headersToSign)
	for _, header := range s.getDefaultHeaders(r) {
		headers = append(headers, header)
	}
	headers = append(headers, [2]string{s.config.AuthHeaderName, authHeader})

	return request.New(
			r.Method(),
			r.RawURL(),
			headers,
			r.Body(),
			r.Expires()),
		nil
}
