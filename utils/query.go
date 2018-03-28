package utils

import (
	"net/url"
	"sort"
)

// ReOrderQuery orders url query for
// hash consistency
func ReOrderQuery(s string) string {
	u, _ := url.Parse(s)
	q, _ := url.ParseQuery(u.RawQuery)

	if len(q) < 1 {
		return s
	}

	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	s = u.Scheme + "://" + u.Host + u.Path + "?"
	for i := 0; i < len(keys); i++ {
		if len(q[keys[i]]) > 0 {
			s += keys[i] + "=" + q[keys[i]][0]
		}
		if i < len(keys)-1 {
			s += "&"
		}
	}

	return s
}
