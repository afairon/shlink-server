package utils

import (
	"net/url"
)

// domains on blacklist
var domains = map[string]struct{}{
	"shlink.cc":   struct{}{},
	"goo.gl":      struct{}{},
	"bit.ly":      struct{}{},
	"tinyurl.com": struct{}{},
	"tiny.cc":     struct{}{},
	"bc.vc":       struct{}{},
}

// IsBlackList checks wether the url is on the
// blacklist or not.
func IsBlackList(s string) (bool, error) {
	u, err := url.Parse(s)
	if err != nil {
		return false, err
	}

	if _, ok := domains[u.Host]; ok {
		return true, nil
	}

	return false, nil
}
