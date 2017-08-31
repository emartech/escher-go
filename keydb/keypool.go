package keydb

import (
	"encoding/json"
	"regexp"
	"sort"
	"strconv"
)

type keyObject struct {
	RawKeyID   string `json:"keyId"`
	Secret     string `json:"secret"`
	AcceptOnly int    `json:"acceptOnly"`
}

type byVersion []keyObject

func (a byVersion) Len() int           { return len(a) }
func (a byVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byVersion) Less(i, j int) bool { return a[i].Version() < a[j].Version() }

var versionMatcher = regexp.MustCompile(`\d+$`)

func (ko keyObject) KeyID() string {
	// if versionMatcher.MatchString(ko.RawKeyID) {
	// 	versionSuffix := fmt.Sprintf("_v%v", ko.Version())

	// 	return strings.TrimSuffix(ko.RawKeyID, versionSuffix)
	// }

	return ko.RawKeyID
}

func (ko keyObject) Version() int {
	if versionMatcher.MatchString(ko.RawKeyID) {
		versionString := versionMatcher.FindString(ko.RawKeyID)
		versionInt, _ := strconv.Atoi(versionString)
		return versionInt
	}

	return 0
}

func parseFromKeyPool(jsonString string) (KeyDB, error) {
	data := []byte(jsonString)

	keypool := make([]keyObject, 0)

	err := json.Unmarshal(data, &keypool)

	if err != nil {
		return nil, err
	}

	keyDB := make(keydb)
	sort.Sort(byVersion(keypool))

	for _, keyObject := range keypool {
		keyID := keyObject.KeyID()

		keyDB[keyID] = append(keyDB[keyID], keyObject.Secret)
	}

	return keyDB, nil
}
