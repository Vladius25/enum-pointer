# EnumPointer

Epointer is a tool to automate the creation of method Pointer() for enums.

For example, given this snippet,

```go
package painkiller

type Pill int

const (
    Placebo Pill = iota
    Aspirin
    Ibuprofen
    Paracetamol
    Acetaminophen = Paracetamol
)
```

running this command

```bash
epointer painkiller_dir
```

in the same directory will create the file epointer_gen.go, in package painkiller,
containing a definition of

```go
func (Pill) Pointer() *Pill
```

That method will return the address of the receiver, so that the call
painkiller.Aspirin.Pointer() will return the address of the Aspirin constant.

Typically this process would be run using go generate, like this:

```bash
go:generate epointer
```

With no arguments, it processes the package in the current directory.
Otherwise, the arguments must name a single directory holding a Go package
or a set of Go source files that represent a single Go package.
