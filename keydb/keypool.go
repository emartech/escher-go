package keydb

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type keyObject struct {
	RawKeyID   string `json:"keyId"`
	Secret     string `json:"secret"`
	AcceptOnly int    `json:"acceptOnly"`
}

var versionMatcher = regexp.MustCompile(`\d+$`)

func (ko keyObject) KeyID() string {
	if versionMatcher.MatchString(ko.RawKeyID) {
		versionSuffix := fmt.Sprintf("_v%v", ko.VersionString())

		return strings.TrimSuffix(ko.RawKeyID, versionSuffix)
	}

	return ko.RawKeyID
}

func (ko keyObject) VersionString() int {
	if versionMatcher.MatchString(ko.RawKeyID) {
		versionInt, _ := strconv.Atoi(versionMatcher.FindString(ko.RawKeyID))

		return versionInt
	}

	return 0
}

func parseFromKeyPool(jsonString string) (KeyDB, error) {
	data := []byte(jsonString)

	keyPool := make([]keyObject, 0)

	err := json.Unmarshal(data, &keyPool)

	if err != nil {
		return nil, err
	}

	keyDB := make(keydb)
	for _, keyObject := range keyPool {
		keyID := keyObject.KeyID()

		keyDB[keyID] = append(keyDB[keyID], keyObject.Secret)
	}

	return keyDB, nil
}
