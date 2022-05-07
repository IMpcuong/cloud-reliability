package main

import (
	"reflect"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	input := []byte("abcdef-12345")
	actual := string(base58Encode(input))
	expected := "2qb7RmPbQXRfszbtQ"
	if reflect.DeepEqual(actual, expected) {
		Error.Fatal("Failed!")
	}
}

func TestBase58Decode(t *testing.T) {
	input := []byte("2qb7RmPbQXRfszbtQ")
	actual := string(base58Decode(input))
	expected := "abcdef-12345"
	if reflect.DeepEqual(actual, expected) {
		Error.Fatal("Failed!")
	}
}
