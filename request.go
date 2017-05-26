package escher

type RequestHeaders [][2]string

type Request struct {
	Method  string         `json:"method"`
	Url     string         `json:"url"`
	Headers RequestHeaders `json:"headers"`
	Body    string         `json:"body"`
}
