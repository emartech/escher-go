package validator

import "errors"

var (
	MissingDateParam           = errors.New("missing date param")
	InvalidEscherKey           = errors.New("Invalid Escher key")
	AlgorithmNotAllowed        = errors.New("Only SHA256 and SHA512 hash algorithms are allowed")
	RequestMethodIsInvalid     = errors.New("The request method is invalid")
	POSTRequestBodyIsEmpty     = errors.New("The request body shouldn't be empty if the request method is POST")
	CredentialDateNotMatching  = errors.New("The credential date does not match with the request date")
	RequestDateNotAcceptable   = errors.New("The request date is not within the accepted time range")
	AuthorizationDateIsInvalid = errors.New("Invalid date in authorization header, it should equal with date header")
	InvalidCredentialScope     = errors.New("The credential scope is invalid")
	HostHeaderNotSigned        = errors.New("The host header is not signed")
	DateHeaderIsNotSigned      = errors.New("The date header is not signed")
	SignatureDoesNotMatch      = errors.New("The signatures do not match")
	MalformedAuthHeader        = errors.New("Could not parse auth header")
	SigningParamNotFound       = errors.New("Signing Param not found")
	SchemaInURLNotAllowed      = errors.New("The request url shouldn't contains http or https")
	HTTPSchemaFoundInTheURL    = errors.New("The request url shouldn't contains http or https")
)
