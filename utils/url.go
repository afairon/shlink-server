package utils

import (
	"net/url"
	"strings"
)

// IsURL check if the string is an URL.
func IsURL(s string) (bool, error) {
	if !strings.HasPrefix(s, "http") {
		return false, nil
	}

	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false, err
	}

	return true, nil
}
