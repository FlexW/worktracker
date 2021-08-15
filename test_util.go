package main

import "testing"

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func assertNoError(t * testing.T, err error) {
	if err != nil {
		t.Fatalf("Error %s", err)
	}
}