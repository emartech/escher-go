package keydb_test

import (
	"testing"

	"github.com/EscherAuth/escher/keydb"
	. "github.com/EscherAuth/escher/testing/env"
	"github.com/stretchr/testify/assert"
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

func TestNewFromENV_KeyPoolValueIsEmpty(t *testing.T) {
	defer UnsetEnvForTheTest(t, "KEY_POOL")()
	keyDB, err := keydb.NewFromENV()
	assert.EqualError(t, err, "KEY_POOL Env value is empty")
	assert.Nil(t, keyDB)
}

func TestNewFromENV_KeyPoolHasKeyObjectAndItIsVersioned(t *testing.T) {
	defer SetEnvForTheTest(t, "KEY_POOL", `[{"keyId":"dpp_ps_v1","secret":"sickrat","acceptOnly":0}]`)()
	keyDB, err := keydb.NewFromENV()

	if err != nil {
		t.Fatal(err)
	}

	secret, err := keyDB.GetSecret("dpp_ps")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "sickrat", secret)

}

// [{"keyId":"dpp_ps_v1","secret":"sickrat","acceptOnly":0}]
