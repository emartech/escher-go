package testing

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/EscherAuth/escher/config"
	"github.com/EscherAuth/escher/keydb"
	"github.com/EscherAuth/escher/request"
)

type TestConfig struct {
	ID                     string
	HeadersToSign          []string           `json:"headersToSign"`
	Title                  string             `json:"title"`
	Description            string             `json:"description"`
	Request                request.Request    `json:"request"`
	Config                 config.Config      `json:"config"`
	Expected               TestConfigExpected `json:"expected"`
	RawKeyDB               [][2]string        `json:"keyDb"`
	FilePath               string
	MandatorySignedHeaders []string `json:"mandatorySignedHeaders"`
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

func fixedConfigBy(tb testing.TB, c config.Config) config.Config {
	var t, err = time.Parse("2006-01-02T15:04:05.999999Z", c.Date)
	if err != nil {
		t, err = time.Parse("Fri, 02 Jan 2006 15:04:05 GMT", c.Date)
	}
	if err != nil {
		t = time.Now().UTC()
	}
	c.Date = t.Format("20060102T150405Z")

	return c
}

func testConfigBy(t testing.TB, filePath string) TestConfig {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		t.Fatal(err)
	}

	var testConfig TestConfig
	json.Unmarshal(content, &testConfig)
	testConfig.FilePath = filePath

	return testConfig
}
