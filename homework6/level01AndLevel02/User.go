package main

const (
	WOMAN Gender = iota
	MAN
)

type Gender uint

type User struct {
	Name    string
	Age     int
	Gender  Gender
	NpyName string
}
