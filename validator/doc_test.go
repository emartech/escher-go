package validator_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
	"github.com/EscherAuth/escher/validator"
)

func ExampleValidator() {

	Config, err := config.NewFromENV()

	if err != nil {
		log.Fatal(err)
	}

	keyDB, err := keydb.NewFromENV()

	if err != nil {
		log.Fatal(err)
	}

	Validator := validator.New(Config)

	handler := func(w http.ResponseWriter, r *http.Request) {

		escherRequest, err := request.NewFromHTTPRequest(r)

		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		apiKey, err := Validator.Validate(escherRequest, keyDB, nil)

		if err != nil {
			w.Header().Set("WWW-Authenticate", "EscherAuth")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		log.Printf("request received from: %v\n", apiKey)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		fmt.Fprintln(w, "OK")

	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler)

	err = http.ListenAndServe(":9292", mux)

	if err != nil {
		log.Fatal(err)
	}

}
