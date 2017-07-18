package keydb_test

import (
	"testing"

	"github.com/EscherAuth/escher/keydb"
	"github.com/stretchr/testify/assert"
)

func TestSecretNotFound(t *testing.T) {

	subject := keydb.NewByKeyValuePair("FOO", "baz")

	_, err := subject.GetSecret("Baz")

	if err == nil {
		t.Fatal("Expected error but nothing raised")
	} else {
		assert.EqualError(t, err, "KeyID Not Found")
	}

}
