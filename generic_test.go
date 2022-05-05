package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestMinimal(t *testing.T) {
	isMinNum := minVal(6, 7)
	isMinStr := minVal("something", "something1")

	if isMinNum != 6 {
		t.Errorf("Failed!")
	}
	if isMinStr != "something" {
		t.Errorf("Failed!")
	}
}

func TestSliceContains(t *testing.T) {
	str := "config/node1/config.json"
	sliceStr := strings.Split(str, "/")
	if !contains(sliceStr, "config.json") {
		t.Errorf("Slice contains unexpected behavior!")
	}
}

func TestSliceRemove(t *testing.T) {
	slice := []int64{1, 2, 3, 5}
	remain := remove(slice, 1)
	if !(len(remain) < len(slice)) {
		t.Errorf("Slice remove length mismatch!")
	}
	// Maybe I will implement `slice.Equal` later!
	if !reflect.DeepEqual(remain, []int64{2, 3, 5}) {
		t.Errorf("Slice remove element mismatch!")
	}
	fmt.Println("The remain slice: ", remain)
}

func TestUniqueSlice(t *testing.T) {
	unq := unique([]int64{1, 1, 2, 2, 3, 5})
	if !reflect.DeepEqual(unq, []int64{1, 2, 3, 5}) {
		t.Errorf("Slice is not unique yet!")
	}
	fmt.Println("Unique slice: ", unq)
}
