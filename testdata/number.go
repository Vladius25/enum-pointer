package main

import "fmt"

type Number int

const (
	_ Number = iota
	One
	Two
	Three
	AnotherOne = One
)

func main() {
	ck(One)
	ck(Two)
	ck(Three)
	ck(AnotherOne)
	ck(127)
}

func ck(num Number) {
	if *num.Pointer() != num {
		panic(fmt.Sprint("number.go: ", *num.Pointer()))
	}
}
