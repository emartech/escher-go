package config

import (
	"errors"
	"strings"

	"github.com/EscherAuth/escher/request"
)

var (
	RequestMethodIsInvalidError  = errors.New("The request method is invalid")
	HTTPSchemaIncludedError      = errors.New("The request url shouldn't contains http or https")
	EmptyRequestBodyForPostError = errors.New("The request body shouldn't be empty if the request method is POST")
	EscherKeyIsInvalidError      = errors.New("Invalid Escher key")
)

func (c Config) Validate(r request.Interface) error {
	if !isMethodAccepted(r) {
		return RequestMethodIsInvalidError
	}

	if isHTTPSchemaPresent(r) {
		return HTTPSchemaIncludedError
	}

	// TODO: Talk with "mr I" about the compatibility issue with AWS or emarsys cases
	// if isPostRequestBodyEmpty(r) {
	// 	return EmptyRequestBodyForPostError
	// }

	if isEscherKeyInvalid(c) {
		return EscherKeyIsInvalidError
	}

	return nil
}

var acceptedMethods = map[string]struct{}{
	"GET":     struct{}{},
	"POST":    struct{}{},
	"PUT":     struct{}{},
	"PATCH":   struct{}{},
	"DELETE":  struct{}{},
	"HEAD":    struct{}{},
	"OPTIONS": struct{}{},
	"LINK":    struct{}{},
	"UNLINK":  struct{}{},
	"TRACE":   struct{}{},
}

func isMethodAccepted(r request.Interface) bool {
	_, ok := acceptedMethods[strings.ToUpper(r.Method())]

	return ok
}

func isHTTPSchemaPresent(r request.Interface) bool {
	u, err := r.URL()
	if err != nil {
		// ignore this case, not this method level scope
		return false
	}

	return u.Scheme != "" && strings.HasPrefix(u.Scheme, "http")
}

func isPostRequestBodyEmpty(r request.Interface) bool {
	if strings.ToLower(r.Method()) != "post" {
		return false
	}

	return r.Body() == ""
}

func isEscherKeyInvalid(c Config) bool {

	if c.ApiSecret == "" || c.AccessKeyId == "" {
		return true
	}

	return false
}
