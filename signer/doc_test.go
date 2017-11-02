package signer_test

import (
	"log"
	"net/http"

	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/request"
	"github.com/EscherAuth/escher/signer"
)

func ExampleSigner() error {

	Config, err := config.NewFromENV()

	if err != nil {
		return err
	}

	req, _ := http.NewRequest("GET", "http://example.com/", nil)

	escherRequest, err := request.NewFromHTTPRequest(req)

	if err != nil {
		return err
	}

	signedRequest, err := signer.New(Config).SignRequest(escherRequest, []string{})

	if err != nil {
		return err
	}

	err = signedRequest.UpdateHTTPRequest(req)

	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, _ := client.Do(req)

	log.Println(resp)

	return nil

}
