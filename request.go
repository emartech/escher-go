package escher

import "net/url"

type Request struct {
	Method  string         `json:"method"`
	Url     string         `json:"url"`
	Headers RequestHeaders `json:"headers"`
	Body    string         `json:"body"`
}

func (r Request) Path() (string, error) {
	url, err := url.Parse(r.Url)

	if err != nil {
		return "", err
	}

	return url.Path, err
}

type QueryParts [][2]string

func (r Request) QueryParts() (QueryParts, error) {
	url, err := url.Parse(r.Url)

	if err != nil {
		return nil, err
	}

	queryParts := make(QueryParts, 0)

	for key, values := range url.Query() {
		for _, value := range values {
			queryParts = append(queryParts, [2]string{key, value})
		}
	}

	return queryParts, nil
}

func (qp QueryParts) Without(key string) QueryParts {
	nqp := make(QueryParts, 0, len(qp))
	for _, kv := range qp {
		if kv[0] != key {
			nqp = append(nqp, kv)
		}
	}
	return nqp
}
