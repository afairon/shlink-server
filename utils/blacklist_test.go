package utils

import (
	"fmt"
	"net/url"
	"testing"
)

func TestIsOnBlacklist(t *testing.T) {
	s := "https://goo.gl/"
	if ban, err := IsOnBlackList(s); !ban || err != nil {
		t.Error("Blacklist not working correctly.")
	}

	s = "http://localhost:8080/"
	u, _ := url.Parse(s)
	fmt.Println(u.Hostname())
	if ban, err := IsOnBlackList(s); !ban || err != nil {
		t.Error("Blacklist not working correctly.")
	}
}
