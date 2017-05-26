package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/adamluzsi/escher-go/runner"
)

var port string

func main() {

	targetPort := "9292"

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	incoming, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	checkForError(err, "could not start server")

	client, err := incoming.Accept()
	checkForError(err, "could not accept cleint connection")

	defer client.Close()

	target, err := net.Dial("tcp", targetPort)
	checkForError(err, "could not connect to target")
	defer target.Close()

	runner.New(targetPort, os.Args[1], os.Args[2:], signals)

}

func checkForError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
