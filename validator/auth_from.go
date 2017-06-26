package validator

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/EscherAuth/escher/request"
)

func (v *validator) getAuthPartsFromHeader(authHeader string) (algorithm, apiKeyID, shortDate, credentialScope string, signedHeaders []string, signature string, expires uint64, err error) {
	expr := regexp.QuoteMeta(v.config.GetAlgoPrefix()) + authHeaderRegexpBase
	rgx, err := regexp.Compile(expr)

	if err != nil {
		return
	}

	m, err := rgxNamedMatch(rgx, authHeader)

	if err != nil {
		return
	}

	algorithm = m["algo"]
	apiKeyID = m["apiKeyID"]
	shortDate = m["shortDate"]
	credentialScope = m["credentials"]
	signedHeaders = strings.Split(m["signedHeaders"], ";")
	signature = m["signature"]

	return
}

func (v *validator) getAuthPartsFromQuery(query request.Query) (algorithm, apiKeyID, shortDate, credentialScope string, signedHeaders []string, signature string, expires uint64, err error) {
	rawExpires, err := v.getSigningParam("Expires", query)
	if err != nil {
		return
	}

	expires, err = strconv.ParseUint(rawExpires, 10, 0)
	if err != nil {
		return
	}

	credential, err := v.getSigningParam("Credentials", query)
	if err != nil {
		return
	}
	credentialParts := strings.SplitN(credential, "/", 3)
	apiKeyID, shortDate, credentialScope = credentialParts[0], credentialParts[1], credentialParts[2]

	rawSignedHeaders, err := v.getSigningParam("SignedHeaders", query)
	if err != nil {
		return
	}
	signedHeaders = strings.Split(rawSignedHeaders, ";")

	rawAlgorithm, err := v.getSigningParam("Algorithm", query)
	if err != nil {
		return
	}

	algorithm, err = v.parseAlgo(rawAlgorithm)
	if err != nil {
		return
	}

	signature, err = v.getSigningParam("Signature", query)

	return
}
