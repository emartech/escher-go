package request

func New(
	method string,
	urlString string,
	headers [][2]string,
	body string,
	expires int,
) (*Request, error) {

	// uri, err := url.Parse(urlString)

	return &Request{
		method:  method,
		url:     urlString,
		headers: headers,
		body:    body,
		expires: expires,
		// UniversalResourceLocator: uri,
	}, nil

}
