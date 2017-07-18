package keydb_test

import (
	"testing"

	"github.com/EscherAuth/escher/keydb"
)

func TestNewByKeyValuePair(t *testing.T) {

	subject := keydb.NewByKeyValuePair("FOO", "baz")

	foundSecret, err := subject.GetSecret("FOO")

	if err != nil {
		t.Fatal(err)
	}

	if foundSecret != "baz" {
		t.Fatalf("expected baz, but actually is %v", foundSecret)
	}

}

func TestNewBySlice(t *testing.T) {

	subject := keydb.NewBySlice([][2]string{[2]string{"hello", "world"}})

	foundSecret, err := subject.GetSecret("hello")

	if err != nil {
		t.Fatal(err)
	}

	if foundSecret != "world" {
		t.Fatalf("expected world, but actually is %v", foundSecret)
	}

}

func TestNewFromENV() {
	defer Set

}