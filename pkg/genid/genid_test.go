package genid

import (
	"testing"
)

func TestIntToBase62(t *testing.T) {
	if id := IntToBase62(0); id != "0" {
		t.Errorf("Expect ID to be '0' but received '%s'", id)
	}

	if id := IntToBase62(1); id != "1" {
		t.Errorf("Expect ID to be 'b' but received '%s'", id)
	}

	if id := IntToBase62(62); id != "10" {
		t.Errorf("Expect ID to be '10' but received '%s'", id)
	}

	if id := IntToBase62(256); id != "48" {
		t.Errorf("Expect ID to be '48' but received '%s'", id)
	}
}
