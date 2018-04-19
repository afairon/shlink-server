package genid

import (
	"fmt"
	"testing"
)

func TestIntToBase62(t *testing.T) {
	if id := IntToBase62(0); id != "a" {
		t.Errorf("Expect ID to be 'a' but received '%s'", id)
	}

	if id := IntToBase62(1); id != "b" {
		t.Errorf("Expect ID to be 'b' but received '%s'", id)
	}

	if id := IntToBase62(62); id != "ba" {
		t.Errorf("Expect ID to be 'ba' but received '%s'", id)
	}

	if id := IntToBase62(256); id != "ei" {
		fmt.Printf("ID: %s\n", id)
		t.Errorf("Expect ID to be 'ei' but received '%s'", id)
	}
}
