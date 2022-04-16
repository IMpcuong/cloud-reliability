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
	Block | BlockChain | Message
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

// MinimalVal compares 2 int numbers to determine which one is smaller.
func MinimalVal[T constraints.Ordered](dst, src T) T {
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

// Contains returns true if a slice contains the given value.
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
