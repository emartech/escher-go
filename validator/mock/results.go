package mock

type ValidationResults []ValidationResult

type ValidationResult struct {
	ApiKey string
	Error  error
}
