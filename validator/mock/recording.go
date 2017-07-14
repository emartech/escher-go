package mock

func (m *Mock) AddValidationResult(apiKey string, err error) {
	m.ValidationResults = append(m.ValidationResults, ValidationResult{ApiKey: apiKey, Error: err})
}
