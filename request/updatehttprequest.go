package request

import "net/http"
import "net/url"

func (ereq *Request) UpdateHTTPRequest(req *http.Request) error {
	mergeHTTPHeaders(ereq, req)

	sURL, err := ereq.URL()

	if err != nil {
		return err
	}

	mergeURLPath(sURL, req.URL)

	return nil
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

func mergeURLPath(s *url.URL, d *url.URL) {
	d.Path = s.Path
	d.RawPath = s.RawPath
	d.RawQuery = s.RawQuery
}
