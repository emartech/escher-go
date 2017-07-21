package mock

import "github.com/EscherAuth/escher/request"

func (m *Mock) AddGeneratedSignature(signature string) {
	m.generatedSignatures = append(m.generatedSignatures, signature)
}

func (m *Mock) nextGeneratedSignature() string {
	signature := m.generatedSignatures[0]
	m.generatedSignatures = m.generatedSignatures[1:]
	return signature
}

func (m *Mock) GenerateSignature(r request.Interface, headersToSign []string) string {
	if len(m.generatedSignatures) == 0 {
		return "signature"
	} else {
		return m.nextGeneratedSignature()
	}
}
