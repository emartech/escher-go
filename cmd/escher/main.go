package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/EscherAuth/escher/cmd/escher/proxy"
	"github.com/EscherAuth/escher/cmd/escher/runner"
)

var port string

func main() {
	targetPort := "2222"
	signals := make(chan os.Signal, 1)
	if len(os.Args) < 3 {
		os.Args = append(os.Args, "-p")
		os.Args = append(os.Args, os.Getenv("PORT"))
	}
	cmd := runner.New(targetPort, os.Args[1], os.Args[2:], signals).Run()
	defer cmd.Wait()

	proxy := proxy.New("http://localhost:" + targetPort)
	http.HandleFunc("/", proxy.Handle)
	http.ListenAndServe(":9292", nil)
	signal.Notify(signals, os.Interrupt)

}

func checkForError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func nem_kell() {
	signals := make(chan os.Signal)
	targetPort := "9393"
	incoming, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	checkForError(err, "could not start server")
	cmd := runner.New(targetPort, os.Args[1], os.Args[2:], signals).Run()
	defer cmd.Wait()

	client, err := incoming.Accept()
	checkForError(err, "could not accept cleint connection")

	target, err := net.Dial("tcp", targetPort)
	checkForError(err, "could not connect to target")
	target.Close()
	client.Close()

}
