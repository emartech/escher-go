package escher

type Escher interface {
	// SignRequest(request Request, headersToSign []string) Request
}

type escher struct {
	config Config
	// signer signer.Signer
}
