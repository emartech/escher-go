package mock

func (m *Mock) popValidationResults() (string, error) {

	if len(m.ValidationResults) == 0 {
		return "", nil
	}

	lastIndex := len(m.ValidationResults) - 1
	lastElement, poppedValidationResults := m.ValidationResults[lastIndex], m.ValidationResults[:lastIndex]

	m.ValidationResults = poppedValidationResults
	return lastElement.ApiKey, lastElement.Error

}
