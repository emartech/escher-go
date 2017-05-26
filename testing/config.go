package testing

import (
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
