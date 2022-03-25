package main

import (
	"testing"
)

func TestGreet(t *testing.T) {
	got := Greet("Gopher")
	expected := "Hello, Gopher!\n"

	if got != expected {
		t.Errorf("Did not get expected result. Wanted %q, got: %q\n", expected, got)
	}
}
