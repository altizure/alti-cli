package text

import (
	"fmt"
	"regexp"
)

// Contains tells whether a contains s.
// And return the index of first match. If not found, index is -1.
func Contains(a []string, s string) (int, bool) {
	for i, e := range a {
		if s == e {
			return i, true
		}
	}
	return -1, false
}

// BestMatch best matches s in a.
// If no exact match is found, will match with regex `\w*s\w*`.
// If no regex match, return def.
func BestMatch(a []string, s, def string) string {
	e := regexp.QuoteMeta(s)
	r, _ := regexp.Compile(fmt.Sprintf("(?i)\\w*%s\\w*", e))
	if i, ok := Contains(a, s); ok {
		return a[i]
	}
	for _, x := range a {
		if r.MatchString(x) {
			return x
		}
	}
	return def
}

// SliceToMap turns a slice of string into a map, with given default value.
func SliceToMap(s []string, d bool) map[string]bool {
	ret := make(map[string]bool)
	for _, v := range s {
		ret[v] = d
	}
	return ret
}
