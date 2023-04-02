// This file contains simple golden tests for various examples.

package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Golden represents a test case.
type Golden struct {
	name   string
	input  string // input; the package clause is provided when running the test.
	output string // expected output.
}

var golden = []Golden{
	{"int_enum", intEnumIn, intEnumOut},
	{"str_enum", strEnumIn, strEnumOut},
	{"conv_enum", convEnumIn, convEnumOut},
	{"mixed_enums", mixedEnumsIn, mixedEnumsOut},
}

// Simple test: enumeration of type int.
const intEnumIn = `type Day int
const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)
`

const intEnumOut = `
func (i Day) Pointer() *Day {
	return &i
}
`

const strEnumIn = `
type Status string
const (
	OK Status = "ok"
	NotOK Status = "not ok"
)
`

const strEnumOut = `
func (i Status) Pointer() *Status {
	return &i
}
`

const convEnumIn = `
type Conv int

const (
	Alpha = Conv(alpha)
	Beta  = Conv(beta)
	Gamma = Conv(gamma)
	Delta = Conv(delta)
)
`

const convEnumOut = `
func (i Conv) Pointer() *Conv {
	return &i
}
`

const mixedEnumsIn = `
type Status string

const (
	posStart = iota
	posMiddle
	posEnd
)

type Weekend int

var (
	Saturday Weekend = iota
	Sunday
)

const (
	OK Status = "ok"
	NotOK Status = "not ok"
)
`

const mixedEnumsOut = `
func (i Status) Pointer() *Status {
	return &i
}
`

func TestGolden(t *testing.T) {
	dir := t.TempDir()
	for _, test := range golden {
		g := Generator{}
		input := "package test\n" + test.input
		file := test.name + ".go"
		absFile := filepath.Join(dir, file)
		err := os.WriteFile(absFile, []byte(input), 0644)
		if err != nil {
			t.Error(err)
		}

		g.parsePackage([]string{absFile}, nil)
		g.generate()
		got := string(g.format())
		if got != test.output {
			t.Errorf("%s: got(%d)\n====\n%q====\nexpected(%d)\n====%q", test.name, len(got), got, len(test.output), test.output)
		}
	}
}
