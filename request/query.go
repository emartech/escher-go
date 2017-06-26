package request

type Query [][2]string

func (q Query) IsInclude(expectedKey string) bool {
	for _, keyValuePair := range q {
		if keyValuePair[0] == expectedKey {
			return true
		}
	}
	return false
}
