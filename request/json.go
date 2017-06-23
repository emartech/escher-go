package request

import (
	"encoding/json"
)

type jsonMapper struct {
	URL     string      `json:"url"`
	Body    string      `json:"body"`
	Method  string      `json:"method"`
	Expires int         `json:"expires"`
	Headers [][2]string `json:"headers"`
}

func ParseJSON(data []byte) (Request, error) {
	r := &request{}
	err := mapJSONContentToRequest(r, data)
	return r, err
}

func (r *request) UnmarshalJSON(data []byte) error {
	return mapJSONContentToRequest(r, data)
}

func mapJSONContentToRequest(r *request, data []byte) error {
	var j jsonMapper
	err := json.Unmarshal(data, &j)

	if err != nil {
		return err
	}

	r.url = j.URL
	r.body = j.Body
	r.method = j.Method
	r.expires = j.Expires
	r.headers = j.Headers

	// uri, err := url.Parse(j.URL)
	// if err != nil {
	// 	return err
	// }
	// r.UniversalResourceLocator = uri

	return nil
}
