package validator

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
	"github.com/EscherAuth/escher/signer"
	"github.com/EscherAuth/escher/utils"
)

func (v *validator) Validate(r request.Interface, keyDB keydb.KeyDB, mandatoryHeaders []string) (string, error) {
	requestForSigning := request.New(r.Method(), r.RawURL(), r.Headers(), v.bodyForSignatureGeneration(r), r.Expires())

	var rawDate string
	var expires uint64
	var algorithm, apiKeyID, shortDate, credentialScope, signature string
	var signedHeaders []string

	headers := r.Headers()
	query := r.Query()

	expectedHeaders := []string{"Host"}
	isSigningInQuery := v.config.IsSigningInQuery(r)

	if !isSigningInQuery {
		expectedHeaders = append(expectedHeaders, v.config.GetAuthHeaderName(), v.config.GetDateHeaderName())
	}

	for _, headerKey := range expectedHeaders {
		_, ok := headers.Get(headerKey)
		if !ok {
			return "", errors.New("The " + strings.ToLower(headerKey) + " header is missing")
		}
	}

	if isSigningInQuery {
		var err error
		rawDate, err = v.getSigningParam("Date", query)

		if err != nil {
			return "", MissingDateParam
		}

		algorithm, apiKeyID, shortDate, credentialScope, signedHeaders, signature, expires, err = v.getAuthPartsFromQuery(query)

		if err != nil {
			return "", err
		}

		requestForSigning.DelQueryValueByKey(v.config.QueryKeyFor("Signature"))
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

		var err error
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
		return "", InvalidEscherKey
	}

	if algorithm != "SHA256" && algorithm != "SHA512" {
		return "", AlgorithmNotAllowed
	}

	if !isValidRequestMethod(r.Method()) {
		return "", RequestMethodIsInvalid
	}

	if strings.ToUpper(r.Method()) == "POST" && r.Body() == "" {
		return "", POSTRequestBodyIsEmpty
	}

	// u, err := url.Parse(request.Url)
	// _, err = url.Parse(r.Url)

	// if err != nil {
	// 	return "", err
	// }

	// if u.Scheme != "" {
	// 	return "", SchemaInURLNotAllowed
	// }

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

	if isSigningInQuery && !isSignedHeadersOnlyInclude(signedHeaders, "host") {
		return "", HostHeaderNotSigned
	}

	if !isSigningInQuery && !isSignedHeadersInlcude(signedHeaders, v.config.GetDateHeaderName()) {
		return "", DateHeaderIsNotSigned
	}

	u, err := r.URL()

	if err != nil {
		return "", err
	}

	if u.Scheme != "" {
		return "", HTTPSchemaFoundInTheURL
	}

	s := signer.New(v.config.Reconfig(rawDate, algorithm, credentialScope, apiKeyID, apiSecret))

	expectedSignature := s.GenerateSignature(requestForSigning, signedHeaders)

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

func (v *validator) getSigningParam(key string, query request.Query) (string, error) {
	queryKey := v.config.QueryKeyFor(key)
	for _, part := range query {
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

	return time.Now().UTC()
}

func (v *validator) bodyForSignatureGeneration(r request.Interface) string {
	var body string

	if v.config.IsSigningInQuery(r) {
		body = "UNSIGNED-PAYLOAD"
	} else {
		body = r.Body()
	}

	return body
}
