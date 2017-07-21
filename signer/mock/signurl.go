package mock

type signURLResponse struct {
	URL   string
	Error error
}

func (m *Mock) AddSignURLResponse(url string, err error) {
	newSignURLResponse := signURLResponse{URL: url, Error: err}
	m.signURLResponses = append(m.signURLResponses, newSignURLResponse)
}

func (m *Mock) nextSignURLResponse() (string, error) {
	first := m.signURLResponses[0]
	m.signURLResponses = m.signURLResponses[1:]
	return first.URL, first.Error
}

func (m *Mock) SignedURLBy(httpMethod, urlToSign string, expires int) (string, error) {
	if len(m.signURLResponses) == 0 {
		return urlToSign, nil
	} else {
		return m.nextSignURLResponse()
	}
}
