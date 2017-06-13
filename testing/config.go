package testing

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	escher "github.com/adamluzsi/escher-go"
	"github.com/adamluzsi/escher-go/keydb"
)

type TestConfig struct {
	ID            string
	HeadersToSign []string           `json:"headersToSign"`
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Request       escher.Request     `json:"request"`
	Config        escher.Config      `json:"config"`
	Expected      TestConfigExpected `json:"expected"`
	RawKeyDB      [][2]string        `json:"keyDb"`
}

func (testConfig TestConfig) KeyDB() keydb.KeyDB {
	return keydb.NewBySlice(testConfig.RawKeyDB)
}

func (testConfig TestConfig) getTitle() string {
	var title string
	if testConfig.Title == "" {
		title = testConfig.ID
	} else {
		title = testConfig.Title + " (" + testConfig.ID + ")"
	}
	return title
}

func fixedConfigBy(tb testing.TB, config escher.Config) escher.Config {
	var t, err = time.Parse("2006-01-02T15:04:05.999999Z", config.Date)
	if err != nil {
		t, err = time.Parse("Fri, 02 Jan 2006 15:04:05 GMT", config.Date)
	}
	if err != nil {
		t = time.Now().UTC()
	}
	config.Date = t.Format("20060102T150405Z")

	return config
}

func testConfigBy(t testing.TB, filePath string) TestConfig {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		t.Fatal(err)
	}

	var testConfig TestConfig
	json.Unmarshal(content, &testConfig)

	return testConfig
}
