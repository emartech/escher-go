package escher

type authenticationHeader struct {
	c Config
	r Request
}

func (r Request) AuthHeaderBy(c Config) *authenticationHeader {
	return &authenticationHeader{c: c, r: r}
}

func (a *authenticationHeader) fetchAuthHeaderString() string {
	for _, headerPair := range a.r.Headers {
		if headerPair[0] == a.c.AuthHeaderName {
			return headerPair[1]
		}
	}
	return ""
}
