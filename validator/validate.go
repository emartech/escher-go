package validator

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	escher "github.com/adamluzsi/escher-go"
	"github.com/adamluzsi/escher-go/keydb"
	"github.com/adamluzsi/escher-go/signer"
	"github.com/adamluzsi/escher-go/utils"
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

	signatureInQuery := v.config.IsSignatureInQuery(request)

	if !signatureInQuery {
		expectedHeaders = append(expectedHeaders, v.config.GetAuthHeaderName(), v.config.GetDateHeaderName())
	}

	for _, headerKey := range expectedHeaders {
		_, ok := request.Headers.Get(headerKey)
		if !ok {
			return "", errors.New("The " + strings.ToLower(headerKey) + " header is missing")
		}
	}

	if method == "GET" && signatureInQuery {
		body = "UNSIGNED-PAYLOAD"

		rawDate, err = v.getSigningParam("Date", queryParts)

		if err != nil {
			return "", MissingDateParam
		}

		algorithm, apiKeyID, shortDate, credentialScope, signedHeaders, signature, expires, err = v.getAuthPartsFromQuery(queryParts)
		if err != nil {
			return "", err
		}

		queryParts = queryParts.Without(v.queryKeyFor("Signature"))
	} else {

		var ok bool
		rawDate, ok = headers.Get(v.config.GetDateHeaderName())

		if !ok {
			return "", errors.New("The " + v.config.GetDateHeaderName() + " header is missing")
		}

		authHeader, ok := headers.Get(v.config.GetAuthHeaderName())
		if !ok {
			return "", errors.New("The " + v.config.GetAuthHeaderName() + " header is missing")
		}

		algorithm, apiKeyID, shortDate, credentialScope, signedHeaders, signature, expires, err = v.getAuthPartsFromHeader(authHeader)
		if err != nil {
			return "", err
		}
	}

	date, err := utils.ParseTime(rawDate)

	if err != nil {
		return "", err
	}

	apiSecret, err := keyDB.GetSecret(apiKeyID)

	if err != nil {
		return "", errors.New("Invalid API key")
	}

	if algorithm != "SHA256" && algorithm != "SHA512" {
		return "", errors.New("Only SHA256 and SHA512 hash algorithms are allowed")
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

	if date.Format("20060102") != shortDate {
		return "", errors.New("The credential date does not match with the request date")
	}

	if !v.isDateWithinRange(date, expires) {
		return "", errors.New("The request date is not within the accepted time range")
	}

	if v.config.ShortDate() != shortDate {
		return "", errors.New("Invalid date in authorization header, it should equal with date header")
	}

	if v.config.CredentialScope != credentialScope {
		return "", errors.New("The credential scope is invalid")
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

	if !signatureInQuery && !isSignedHeadersInlcude(signedHeaders, v.config.GetDateHeaderName()) {
		return "", errors.New("The date header is not signed")
	}

	s := signer.New(v.config.Reconfig(date.Format(utils.EscherDateFormat), algorithm, credentialScope, apiKeyID, apiSecret))

	expectedSignature := s.GenerateSignature(request, signedHeaders)

	if expectedSignature != signature {
		return "", errors.New("The signatures do not match")
	}

	return apiKeyID, nil
}

func isSignedHeadersOnlyInclude(signedHeaders []string, keyword string) bool {
	return isSignedHeadersInlcude(signedHeaders, keyword) && len(signedHeaders) == 1
}

const authHeaderRegexpBase = "-HMAC-(?P<algo>[A-Z0-9\\,]+) Credential=(?P<apiKeyID>[A-Za-z0-9\\-_]+)/(?P<shortDate>[0-9]{8})/(?P<credentials>[A-Za-z0-9\\-_ /]+), SignedHeaders=(?P<signedHeaders>[A-Za-z\\-;]+), Signature=(?P<signature>[0-9a-f]+)$"

func rgxNamedMatch(rgx *regexp.Regexp, text string) (map[string]string, error) {

	match := rgx.FindStringSubmatch(text)
	subexpNames := rgx.SubexpNames()

	if len(match) != len(subexpNames) {
		return nil, errors.New("Could not parse auth header")
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
	return "X-" + v.config.GetVendorKey() + "-" + key
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

const parseAlgoRgxBase = "-HMAC-(?P<algo>[A-Z0-9\\,]+)$"

func (v *validator) parseAlgo(algorithm string) (string, error) {
	rgx, err := regexp.Compile("^" + regexp.QuoteMeta(v.config.GetAlgoPrefix()) + parseAlgoRgxBase)
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

func (v *validator) isDateWithinRange(t time.Time, expires uint64) bool {
	timeNow := v.Time()
	min := t.Add(-1 * clockSkew * time.Second)
	max := t.Add(time.Duration(clockSkew+expires) * time.Second)
	return min.Before(timeNow) && max.After(timeNow)
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

func (v *validator) Time() time.Time {
	if v.config.Date != "" {
		t, err := utils.ParseTime(v.config.Date)
		if err != nil {
			panic(err)
		}
		return t
	}

	return time.Now()
}
