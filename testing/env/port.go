package env

import (
	"net"
	"strconv"
	"testing"
)

func FatalIfPortIsAlreadyInUse(t testing.TB, port int) {

	conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))

	if err == nil {
		conn.Close()
		t.Fatal("port shouldn't listen!")
	}

}
