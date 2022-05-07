package main

import (
	"os"
	"sort"
	"strconv"

	"golang.org/x/exp/constraints"
)

// This is the utilities class for generalize/abstract functions with
// akin behaviors can be implemented for multiple use-cases.

type Utils interface {
	Serialize() []byte
	Deserialize(encoded []byte) *any
	Block | Blockchain | Message
}

// ByModTime sorting files with the order of the last modification time.
type ByModTime []os.FileInfo

// Len method get length of all (modification) files.
func (fm ByModTime) Len() int {
	return len(fm)
}

// Swap method swap two given files.
func (fm ByModTime) Swap(src, dst int) {
	fm[src], fm[dst] = fm[dst], fm[src]
}

// Less method to compare two files by modification time.
func (fm ByModTime) Less(src, dst int) bool {
	return fm[src].ModTime().Before(fm[dst].ModTime())
}

// SortFiles sorts the all the files inside the given directory
// with the order of their modification time.
func SortFiles(dir string) {
	files, err := os.Open(dir)
	if err != nil {
		Error.Printf("Could not open directory %s", dir)
	}
	fm, err := files.Readdir(-1)
	if err != nil {
		Error.Printf("Could not read files %v", files)
	}
	defer files.Close()
	sort.Sort(ByModTime(fm))
}

// minVal compares 2 integer numbers to determine which one is smaller.
func minVal[T constraints.Ordered](dst, src T) T {
	if dst < src {
		return dst
	}
	return src
}

// Itobytes converts an int number to a bytes slice.
func Itobytes(val int) []byte {
	return []byte(strconv.Itoa(val))
}

// Bytestoi converts a bytes slice to an int number.
func Bytestoi(val []byte) int {
	i, err := strconv.Atoi(string(val))
	if err != nil {
		Error.Panic(err)
	}
	return i
}

// contains returns true if slice `s` contains the given element `e`.
func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// remove the given element `e` from slice `s`.
func remove[T comparable](s []T, e T) []T {
	if contains(s, e) {
		pos := indexOf(s, e)
		return append(s[:pos], s[pos+1:]...)
	}
	return s
}

// indexOf returns the index of the first occurrence of the provided `e` in `s`.
func indexOf[T comparable](s []T, e T) int {
	for pos, v := range s {
		if e == v {
			return pos
		}
	}
	return -1
}

// unique returns a unique slice with no duplicated values.
func unique[T comparable](s []T) []T {
	unqMap := make(map[T]bool)
	var res []T
	for _, v := range s {
		if _, ok := unqMap[v]; !ok {
			unqMap[v] = true
			res = append(res, v)
		}
	}
	return res
}

// reverseBytes reverse the order of a byte slice.
func reverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
