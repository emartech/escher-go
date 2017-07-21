package mock_test

import (
	"math/rand"
	"strconv"

	"github.com/EscherAuth/escher/request"
)

func requestBy(method, url string) *request.Request {
	return request.New(method, url, [][2]string{}, strconv.Itoa(rand.Intn(42)), rand.Intn(42))
}
