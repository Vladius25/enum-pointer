// Check that constants defined as a conversion are accepted.

package main

import "fmt"

type Other int // Imagine this is in another package.

const (
	alpha Other = iota
	beta
	gamma
	delta
)

type Conv int

const (
	Alpha = Conv(alpha)
	Beta  = Conv(beta)
	Gamma = Conv(gamma)
	Delta = Conv(delta)
)

func main() {
	ck(Alpha)
	ck(Beta)
	ck(Gamma)
	ck(Delta)
	ck(42)
}

func ck(c Conv) {
	if *c.Pointer() != c {
		panic(fmt.Sprint("conv.go: ", *c.Pointer()))
	}
}
