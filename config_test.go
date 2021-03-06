package main

import (
	"fmt"
	"testing"
)

func TestFilePathWalkDir(t *testing.T) {
	actual, err := walkCfgDir("")
	if err != nil {
		t.Errorf("Cannot read directory!")
	}
	if actual != nil {
		fmt.Printf("actual: %v\n", actual)
	}
}

func TestFilePathAnyDir(t *testing.T) {
	actual, err := walkCfgDir(".")
	if err != nil {
		t.Errorf("Cannot read directory!")
	}
	if actual != nil {
		fmt.Printf("actual: %v\n", actual)
	}
}

func TestIOReadFilePath(t *testing.T) {
	expected := []string{"node1", "node2", "node3"}
	files, _ := walkCfgDir("")
	fmt.Printf("=====> %v\n", files)
	for i, file := range files {
		actual, err := readPaths(file)
		if err != nil {
			t.Errorf("Cannot read directory!")
		}
		if actual != nil {
			fmt.Printf("actual: idx=%d paths=%v len=%v\n", i, actual, len(actual))
			fmt.Printf("expected: %v\n", expected)
		}
	}
}
