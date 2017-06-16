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

func (v *validator) Validate(request escher.Request, keyDB keydb.KeyDB, mandatoryHeaders []string) (string, error) {

	requestForSigning := &escher.Request{
		Method:  request.Method,
		Url:     request.Url,
		Headers: request.Headers,
		Body:    request.Body,
	}

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

	if signatureInQuery {
		requestForSigning.Body = "UNSIGNED-PAYLOAD"

		rawDate, err = v.getSigningParam("Date", queryParts)

		if err != nil {
			return "", MissingDateParam
		}

		algorithm, apiKeyID, shortDate, credentialScope, signedHeaders, signature, expires, err = v.getAuthPartsFromQuery(queryParts)

		if err != nil {
			return "", err
		}

		requestForSigning.DelQueryValueByKey(v.config.QueryKeyFor("Signature"))
	} else {

		var ok bool
		rawDate, ok = request.Headers.Get(v.config.GetDateHeaderName())

		if !ok {
			return "", errors.New("The " + v.config.GetDateHeaderName() + " header is missing")
		}

		authHeader, ok := request.Headers.Get(v.config.GetAuthHeaderName())

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
		return "", InvalidAPIKey
	}

	if algorithm != "SHA256" && algorithm != "SHA512" {
		return "", AlgorithmNotAllowed
	}

	if !isValidRequestMethod(request.Method) {
		return "", RequestMethodIsInvalid
	}

	if strings.ToUpper(request.Method) == "POST" && request.Body == "" {
		return "", POSTRequestBodyIsEmpty
	}

	_, err = request.Path()

	if err != nil {
		return "", err
	}

	if date.Format("20060102") != shortDate {
		return "", CredentialDateNotMatching
	}

	if !v.isDateWithinRange(date, expires) {
		return "", RequestDateNotAcceptable
	}

	if v.config.ShortDate() != shortDate {
		return "", AuthorizationDateIsInvalid
	}

	if v.config.CredentialScope != credentialScope {
		return "", InvalidCredentialScope
	}

	if !isSignedHeadersInlcude(signedHeaders, "host") {
		return "", HostHeaderNotSigned
	}

	if mandatoryHeaders != nil {
		for _, headerKey := range mandatoryHeaders {
			if !isSignedHeadersInlcude(signedHeaders, headerKey) {
				return "", errors.New("The " + headerKey + " header is not signed")
			}
		}
	}

	if signatureInQuery && !isSignedHeadersOnlyInclude(signedHeaders, "host") {
		return "", HostHeaderNotSigned
	}

	if !signatureInQuery && !isSignedHeadersInlcude(signedHeaders, v.config.GetDateHeaderName()) {
		return "", DateHeaderIsNotSigned
	}

	s := signer.New(v.config.Reconfig(date.Format(utils.EscherDateFormat), algorithm, credentialScope, apiKeyID, apiSecret))

	expectedSignature := s.GenerateSignature(*requestForSigning, signedHeaders)

	if expectedSignature != signature {
		return "", SignatureDoesNotMatch
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
		return nil, MalformedAuthHeader
	}

	result := make(map[string]string)
	for i, name := range subexpNames {
		if i != 0 {
			result[name] = match[i]
		}
	}

	return result, nil
}

func (v *validator) getSigningParam(key string, queryParts escher.QueryParts) (string, error) {
	queryKey := v.config.QueryKeyFor(key)
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
