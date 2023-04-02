// This file has a build tag "tag"

//go:build tag
// +build tag

package main

type ProtectedConst int

const TagProtected ProtectedConst = C + 1
