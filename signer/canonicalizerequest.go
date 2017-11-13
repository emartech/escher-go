package signer

import (
	"strings"

	"github.com/EscherAuth/escher/request"
)

func (s *signer) CanonicalizeRequest(r request.Interface, headersToSign []string) string {
	var u = parsePathQuery(r.RawURL())
	parts := make([]string, 0, 6)
	parts = append(parts, r.Method())
	parts = append(parts, canonicalizePath(u.Path))
	parts = append(parts, canonicalizeQuery(u.Query))
	parts = append(parts, s.canonicalizeHeaders(r, headersToSign))
	parts = append(parts, s.canonicalizeHeadersToSign(r, headersToSign))
	parts = append(parts, s.computeDigest(r.Body()))
	canonicalizedRequest := strings.Join(parts, "\n")
	return canonicalizedRequest
}
