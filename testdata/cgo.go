// Import "C" shouldn't be imported.

package main

/*
#define HELLO 1
*/
import "C"
import "fmt"

type Cgo uint32

const (
	// MustScanSubDirs indicates that events were coalesced hierarchically.
	MustScanSubDirs Cgo = 1 << iota
)

func main() {
	_ = C.HELLO
	ck(MustScanSubDirs)
}

func ck(cgo Cgo) {
	if *cgo.Pointer() != cgo {
		panic(fmt.Sprintf("cgo.go: %v", *cgo.Pointer()))
	}
}
