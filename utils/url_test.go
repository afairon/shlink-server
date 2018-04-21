package utils

import (
	"testing"
)

func TestIsURL(t *testing.T) {
	url := ""
	if b, ok := IsURL(url); b {
		t.Error(b, ok)
		t.Error("URL identification error. This URL is incorrect.")
	}

	url = "mongodb://localhost:27017/test"
	if b, ok := IsURL(url); b {
		t.Error(b, ok)
		t.Error("URL identification error. This URL is incorrect.")
	}

	url = "https://localhost/"
	if b, ok := IsURL(url); !b || ok != nil {
		t.Error("URL identification error. This URL is correct.")
	}
}
