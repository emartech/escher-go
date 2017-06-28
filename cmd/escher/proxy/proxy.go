package proxy

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
)

type prox struct {
	// target url of reverse proxy
	target *url.URL
	// instance of Go ReverseProxy thatwill do the job for us
	proxy *httputil.ReverseProxy
}

// small factory
func New(target string) *prox {
	url, _ := url.Parse(target)
	// you should handle error on parsing
	return &prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *prox) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")
	validator := Validator{}
	keyDb := keydb.NewBySlice([][2]string{})
	_, err := validator.Validate(request.Request{}, keyDb, nil)
	if err != nil {
		out, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		r.Body = NewMyReader(string(out))
		p.proxy.ServeHTTP(w, r)
	}

}

type myReader struct {
	s *strings.Reader
}

func (m *myReader) Close() error {
	return nil
}

func (m *myReader) Read(p []byte) (n int, err error) {
	return m.s.Read(p)
}

func NewMyReader(s string) *myReader {
	return &myReader{s: strings.NewReader(s)}
}

type Validator struct {
}

func (v *Validator) Validate(r request.Request, keyDB keydb.KeyDB, mandatoryHeaders []string) (string, error) {
	return "API_KEY", errors.New("KACSA")
}
