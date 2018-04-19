package utils

import "testing"

func TestQueryOrder(t *testing.T) {
	url := "http://localhost:8080/add?name=Doe&firstname=John"
	expected := "http://localhost:8080/add?firstname=John&name=Doe"
	if url := ReOrderQuery(url); url != expected {
		t.Error("Query not sorted correctly.")
	}
}
