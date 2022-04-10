package main

import (
	"fmt"
	"testing"
)

func TestIORead(t *testing.T) {
	expected := []string{"node1", "node2", "node3"}
	actual, err := IOReadDir("")
	if err != nil {
		t.Errorf("Cannot read directory!")
	}
	if actual != nil {
		fmt.Printf("actual: %v\n", actual)
		fmt.Printf("expected: %v\n", expected)
	}
}
