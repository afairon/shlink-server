package utils

import (
	"testing"
)

func TestIsBlacklist(t *testing.T) {
	url := "https://goo.gl/"
	if ban, ok := IsBlackList(url); !ban || ok != nil {
		t.Error("Blacklist not working correctly.")
	}
}
