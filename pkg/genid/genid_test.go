package genid

import (
	"fmt"
	"testing"
)

func TestGenerateNextID(t *testing.T) {
	code, err := GenerateNextID("")
	if err != nil {
		t.Error(err)
	}

	if code != "a" {
		fmt.Printf("ID: %s\n", code)
		t.Error("Expect ID to be a")
	}

	code, err = GenerateNextID("a")
	if err != nil {
		t.Error(err)
	}

	if code != "b" {
		fmt.Printf("ID: %s\n", code)
		t.Error("Expect ID to be b")
	}

	code, err = GenerateNextID("9")
	if err != nil {
		t.Error(err)
	}

	if code != "aa" {
		fmt.Printf("ID: %s\n", code)
		t.Error("Expect ID to be aa")
	}

	code, err = GenerateNextID("aa")
	if err != nil {
		t.Error(err)
	}

	if code != "ab" {
		fmt.Printf("ID: %s\n", code)
		t.Error("Expect ID to be ab")
	}

	code, err = GenerateNextID("a9")
	if err != nil {
		t.Error(err)
	}

	if code != "ba" {
		fmt.Printf("ID: %s\n", code)
		t.Error("Expect ID to be ba")
	}
}

func BenchmarkGenID(b *testing.B) {
	var id string
	for i := 0; i <= 1000; i++ {
		code, err := GenerateNextID(id)
		if err != nil {
			b.Fatal("Error detected")
		}
		id = code
	}
}
