package utils

import (
	"testing"
)

func TestIsBlacklist(t *testing.T) {
	url := "https://goo.gl/"
	if ban, err := IsBlackList(url); !ban || err != nil {
		t.Error("Blacklist not working correctly.")
	}
}
