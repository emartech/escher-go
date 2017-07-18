package env

import (
	"os"
	"testing"
)

func SetEnvForTheTest(t testing.TB, key, value string) func() {
	restorer := envRestorerFor(t, key)

	err := os.Setenv(key, value)

	if err != nil {
		t.Fatal(err)
	}

	return restorer
}

func UnsetEnvForTheTest(t testing.TB, key string) func() {
	restorer := envRestorerFor(t, key)

	err := os.Unsetenv(key)

	if err != nil {
		t.Fatal(err)
	}

	return restorer
}

func envRestorerFor(t testing.TB, key string) func() {

	originalValue, keyWasSet := os.LookupEnv(key)

	return func() {

		var err error

		if keyWasSet {
			err = os.Setenv(key, originalValue)
		} else {
			err = os.Unsetenv(key)
		}

		if err != nil {
			t.Fatal(err)
		}

	}

}
