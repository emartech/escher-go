package mock

type Mock struct {
	signRequestResponses []signRequestResponse
	signURLResponses     []signURLResponse
	generatedSignatures  []string 
}

func New() *Mock {
	return &Mock{}
}
