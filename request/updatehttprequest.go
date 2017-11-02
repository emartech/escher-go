package request

import "net/http"

func (ereq *Request) UpdateHTTPRequest(req *http.Request) error {
	mergeHTTPHeaders(ereq, req)
	return mergeHTTPURL(ereq, req)
}

func mergeHTTPHeaders(s *Request, d *http.Request) {
	Header := make(http.Header)

	for _, KeyValue := range s.Headers() {
		Header.Add(KeyValue[0], KeyValue[1])
	}

	for name, values := range d.Header {
		if _, ok := Header[name]; !ok {
			Header[name] = values
		}
	}

	d.Header = Header
}

func mergeHTTPURL(s *Request, d *http.Request) error {
	sURL, err := s.URL()

	if err != nil {
		return err
	}

	d.URL = sURL

	return nil
}
