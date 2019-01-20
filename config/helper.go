package config

import (
	"fmt"
	"net/url"
	"strings"
)

// endpointToKey returns the unique key of this scope.
// The scope is defined by the scheme, host and port of endpoint.
func endpointToKey(ep string) string {
	u, err := url.ParseRequestURI(ep)
	if err != nil {
		return strings.Replace(DefaultScope, ".", "*", -1)
	}
	str := fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
	if u.Port() != "" {
		str += ":" + u.Port()
	}
	return strings.Replace(str, ".", "*", -1)
}

// uniqueProfile returns a new unique set of profiles.
func uniqueProfile(ps []Profile) []Profile {
	var ret []Profile
	for _, x := range ps {
		found := false
		for _, y := range ret {
			if x.Equal(y) {
				found = true
				break
			}
		}
		if !found {
			ret = append(ret, x)
		}
	}
	return ret
}

// compareStrs computes the number of same runes in the prefix of both strings.
func compareStrs(a, b string) int {
	ret := 0
	m := min(len(a), len(b))
	for i := 0; i < m; i++ {
		if a[i] == b[i] {
			ret++
		} else {
			break
		}
	}
	return ret
}

// min returns the minimum of two.
func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}
