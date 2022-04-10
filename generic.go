package main

import "strconv"

// This is the utilities class for generalize/abstract functions with
// akin behaviors can be implemented for multiple use-cases.

type Utils interface {
	Serialize() []byte
	Deserialize(encoded []byte) *any
	Block | BlockChain | Message
}

// MinimalVal compares 2 int numbers to determine which one is smaller.
func MinimalVal(dst, src int) int {
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
