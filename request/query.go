package request

type Query [][2]string

func (r *Request) Query() Query {
	return r.parsePathQuery().Query
}

func (r *Request) DelQueryValueByKey(key string) error {
	uri, err := r.URL()

	if err != nil {
		return err
	}

	values := uri.Query()
	values.Del(key)
	uri.RawQuery = values.Encode()

	r.url = uri.String()

	return nil
}

func (q Query) IsInclude(expectedKey string) bool {
	for _, keyValuePair := range q {
		if keyValuePair[0] == expectedKey {
			return true
		}
	}
	return false
}
