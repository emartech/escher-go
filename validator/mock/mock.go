package mock

type Mock struct {
	ReceivedValidations []ReceivedValidation
	ValidationResults   []ValidationResult
}

func New() *Mock { return &Mock{} }
