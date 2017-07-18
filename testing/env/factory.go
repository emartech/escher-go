package env

import (
	"os"
	"testing"
)

func SetEnvForTheTest(t testing.TB, key, value string) func() {

	orgEnvValue, envKeyWasSetBeforeTheTest := os.LookupEnv(key)

	err := os.Setenv(key, value)

	if err != nil {
		t.Fatal(err)
	}

	return func() {

		var err error
		if envKeyWasSetBeforeTheTest {
			err = os.Setenv(key, orgEnvValue)
		} else {
			err = os.Unsetenv(key)
		}

		if err != nil {
			t.Fatal(err)
		}

	}

}
