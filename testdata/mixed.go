package main

import "fmt"

type (
	Status string
	Unum   uint8
)

const (
	StatusOK    Status = "OK"
	m2          Unum   = 1
	StatusError Status = "ERROR"
	m1          Unum   = 2
)

func main() {
	ckStatus(StatusOK)
	ckStatus(StatusError)
	ckUnum(m1)
	ckUnum(m2)
}

func ckStatus(status Status) {
	if *status.Pointer() != status {
		panic(fmt.Sprint("mixed.go: ", *status.Pointer()))
	}
}

func ckUnum(unum Unum) {
	if *unum.Pointer() != unum {
		panic(fmt.Sprint("mixed.go: ", *unum.Pointer()))
	}
}
