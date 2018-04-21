package utils

import (
	"os"
	"testing"
)

func TestReadConfig1(t *testing.T) {
	config := New()

	if ok, _ := exists("shlink.yml"); !ok {
		config.ReadConfig()
	}
}

func TestReadConfig2(t *testing.T) {
	config := New()

	if ok, err := exists("shlink.yml"); !ok || err != nil {
		t.Fatalf("Err: %s", err)
	}

	config.ReadConfig()
	os.Remove("shlink.yml")
}
