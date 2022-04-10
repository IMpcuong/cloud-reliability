package main

// This is the utilities class for generalize/abstract functions with
// akin behaviors can be implemented for multiple use-cases.

type Utils interface {
	Serialize() []byte
	Deserialize(encoded []byte) *any
}
