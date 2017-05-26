package validator

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	escher "github.com/adamluzsi/escher-go"
	"github.com/adamluzsi/escher-go/keydb"
	"github.com/adamluzsi/escher-go/signer"
)

var (
	MissingDateParam = errors.New("missing date param")
)

func (v *validator) Validate(request escher.Request, keyDB keydb.KeyDB, mandatoryHeaders []string) (string, error) {

	method := request.Method
	body := request.Body
	headers := request.Headers

	var rawDate string
	var expires uint64
	var algorithm, apiKeyID, shortDate, credentialScope, signature string
	var signedHeaders []string

	queryParts, err := request.QueryParts()
	if err != nil {
		return "", err
	}

	expectedHeaders := []string{"Host"}

	_, signatureIsNotInQuery := v.getSigningParam("Signature", queryParts)
	signatureInQuery := signatureIsNotInQuery == nil

	if signatureIsNotInQuery != nil {
		expectedHeaders = append(expectedHeaders, v.config.AuthHeaderName, v.config.DateHeaderName)
	}

	for _, headerKey := range expectedHeaders {
		_, ok := request.Headers.Get(headerKey)
		if !ok {
			return "", errors.New("The " + strings.ToLower(headerKey) + " header is missing")
		}
	}

	if method == "GET" && signatureIsNotInQuery == nil {
		body = "UNSIGNED-PAYLOAD"

		rawDate, err := v.getSigningParam("Date", queryParts)
		_ = rawDate

		if err != nil {
			return "", MissingDateParam
		}

		algorithm, apiKeyID, shortDate, credentialScope, signedHeaders, signature, expires, err = v.getAuthPartsFromQuery(queryParts)
		if err != nil {
			return "", err
		}

		// queryParts.delete [query_key_for('Signature'), signature]
		// queryParts = queryParts.map { |k, v| [k, v] }
	} else {

		rawDate, ok := headers.Get(v.config.DateHeaderName)
		_ = rawDate

		if !ok {
			return "", errors.New("The " + v.config.DateHeaderName + " header is missing")
		}

		authHeader, ok := headers.Get(v.config.AuthHeaderName)
		if !ok {
			return "", errors.New("The " + v.config.AuthHeaderName + " header is missing")
		}

		algorithm, apiKeyID, shortDate, credentialScope, signedHeaders, signature, expires, err = v.getAuthPartsFromHeader(authHeader)
		if err != nil {
			return "", err
		}
	}

	date, err := time.Parse(time.RFC3339, rawDate)

	if err != nil {
		return "", err
	}

	apiSecret, err := keyDB.GetSecret(apiKeyID)

	if err != nil {
		return "", errors.New("Invalid Escher key")
	}

	if algorithm != "SHA256" && algorithm != "SHA512" {
		return "", errors.New("Invalid hash algorithm, only SHA256 and SHA512 are allowed")
	}

	if !isValidRequestMethod(method) {
		return "", errors.New("The request method is invalid")
	}

	if strings.ToUpper(method) == "POST" && body == "" {
		return "", errors.New("The request body shouldn't be empty if the request method is POST")
	}

	// raise EscherError, "The request url shouldn't contains http or https" if path.match /^https?:\/\//
	_, err = request.Path()

	if err != nil {
		return "", err
	}

	if v.config.ShortDate() != shortDate {
		return "", errors.New("Invalid date in authorization header, it should equal with date header")
	}

	if isDateWithinRange(date, expires) {
		return "", errors.New("The request date is not within the accepted time range")
	}

	if v.config.CredentialScope != credentialScope {
		return "", errors.New("Invalid Credential Scope")
	}

	if !isSignedHeadersInlcude(signedHeaders, "host") {
		return "", errors.New("The host header is not signed")
	}

	if mandatoryHeaders != nil {
		for _, headerKey := range mandatoryHeaders {
			if !isSignedHeadersInlcude(signedHeaders, headerKey) {
				return "", errors.New("The " + headerKey + " header is not signed")
			}
		}
	}

	if signatureInQuery && !isSignedHeadersOnlyInclude(signedHeaders, "host") {
		return "", errors.New("Only the host header should be signed")
	}

	if !signatureInQuery && !isSignedHeadersInlcude(signedHeaders, v.config.DateHeaderName) {
		return "", errors.New("The date header is not signed")
	}

	s := signer.New(v.config.Reconfig(date, apiKeyID, apiSecret))
	expectedSignature := s.GenerateSignature(request, signedHeaders)

	if expectedSignature != signature {
		return "", errors.New("The signatures do not match")
	}

	return apiKeyID, nil
}

func isSignedHeadersOnlyInclude(signedHeaders []string, keyword string) bool {
	return isSignedHeadersInlcude(signedHeaders, keyword) && len(signedHeaders) == 1
}

const authHeaderRegexpBase = "-HMAC-(?<algo>[A-Z0-9\\,]+) Credential=(?<apiKeyID>[A-Za-z0-9\\-_]+)/(?<shortDate>[0-9]{8})/(?<credentials>[A-Za-z0-9\\-_ /]+), SignedHeaders=(?<signedHeaders>[A-Za-z\\-;]+), Signature=(?<signature>[0-9a-f]+)$"

func (v *validator) getAuthPartsFromHeader(authHeader string) (algorithm, apiKeyID, shortDate, credentialScope string, signedHeaders []string, signature string, expires uint64, err error) {
	expr := regexp.QuoteMeta(v.config.AlgoPrefix) + authHeaderRegexpBase
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
	signedHeaders = strings.Split(m["signedHeaders"], ":")
	signature = m["signature"]

	return
}

func rgxNamedMatch(rgx *regexp.Regexp, text string) (map[string]string, error) {

	match := rgx.FindStringSubmatch(text)
	subexpNames := rgx.SubexpNames()

	if len(match) != len(subexpNames) {
		return nil, errors.New("regexp not matchable")
	}

	result := make(map[string]string)
	for i, name := range subexpNames {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return result, nil
}

func (v *validator) queryKeyFor(key string) string {
	return "X-" + v.config.VendorKey + "-" + key
}

var SigningParamNotFound = errors.New("Signing Param not found")

func (v *validator) getSigningParam(key string, queryParts escher.QueryParts) (string, error) {
	queryKey := v.queryKeyFor(key)
	for _, part := range queryParts {
		if part[0] == queryKey {
			return url.QueryUnescape(part[1])
		}
	}

	return "", SigningParamNotFound
}

func (v *validator) getAuthPartsFromQuery(queryParts escher.QueryParts) (algorithm, apiKeyID, shortDate, credentialScope string, signedHeaders []string, signature string, expires uint64, err error) {
	rawExpires, err := v.getSigningParam("Expires", queryParts)
	if err != nil {
		return
	}

	expires, err = strconv.ParseUint(rawExpires, 10, 0)
	if err != nil {
		return
	}

	credential, err := v.getSigningParam("Credentials", queryParts)
	if err != nil {
		return
	}
	credentialParts := strings.SplitN(credential, "/", 3)
	apiKeyID, shortDate, credentialScope = credentialParts[0], credentialParts[1], credentialParts[2]

	rawSignedHeaders, err := v.getSigningParam("SignedHeaders", queryParts)
	if err != nil {
		return
	}
	signedHeaders = strings.Split(rawSignedHeaders, ";")

	rawAlgorithm, err := v.getSigningParam("Algorithm", queryParts)
	if err != nil {
		return
	}

	algorithm, err = v.parseAlgo(rawAlgorithm)
	if err != nil {
		return
	}

	signature, err = v.getSigningParam("Signature", queryParts)

	return
}

const parseAlgoRgxBase = "-HMAC-(?<algo>[A-Z0-9\\,]+)$"

func (v *validator) parseAlgo(algorithm string) (string, error) {
	rgx, err := regexp.Compile("^" + regexp.QuoteMeta(v.config.AlgoPrefix) + parseAlgoRgxBase)
	if err != nil {
		return "", err
	}

	dictionary, err := rgxNamedMatch(rgx, algorithm)
	if err != nil {
		return "", err
	}

	return dictionary["algo"], nil
}

const clockSkew = 300

func isDateWithinRange(t time.Time, expires uint64) bool {
	timeNow := time.Now()

	return t.Add(-1*clockSkew*time.Second).After(timeNow) && t.Add(time.Duration(clockSkew+expires)*time.Second).Before(timeNow)
}

var acceptedRequestMethods = map[string]struct{}{
	"OPTIONS": struct{}{},
	"GET":     struct{}{},
	"HEAD":    struct{}{},
	"POST":    struct{}{},
	"PUT":     struct{}{},
	"DELETE":  struct{}{},
	"TRACE":   struct{}{},
	"PATCH":   struct{}{},
	"CONNECT": struct{}{},
}

func isValidRequestMethod(method string) bool {
	_, ok := acceptedRequestMethods[strings.ToUpper(method)]
	return ok
}

func isSignedHeadersInlcude(signedHeaders []string, keyword string) bool {
	formattedKeyword := strings.ToLower(keyword)
	for _, headerName := range signedHeaders {
		if strings.ToLower(headerName) == formattedKeyword {
			return true
		}
	}
	return false
}
