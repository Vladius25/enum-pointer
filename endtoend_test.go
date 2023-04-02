package main

import (
	"bytes"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// This file contains a test that compiles and runs each program in testdata
// after generating the pointer method for its type. The rule is that for testdata/x.go
// we run epointer and then compile and run the program. The resulting
// binary panics if the Pointer methods are not correct, including for error cases.

func TestMain(m *testing.M) {
	if os.Getenv("EPOINTER_TEST_IS_EPOINTER") != "" {
		main()
		os.Exit(0)
	}

	// Inform subprocesses that they should run the cmd/epointer main instead of
	// running tests. It's a close approximation to building and running the real
	// command, and much less complicated and expensive to build and clean up.
	os.Setenv("EPOINTER_TEST_IS_EPOINTER", "1")

	os.Exit(m.Run())
}

func TestEndToEnd(t *testing.T) {
	epointer := epointerPath(t)
	// Read the testdata directory.
	fd, err := os.Open("testdata")
	if err != nil {
		t.Fatal(err)
	}
	defer fd.Close()
	names, err := fd.Readdirnames(-1)
	if err != nil {
		t.Fatalf("Readdirnames: %s", err)
	}
	// Generate, compile, and run the test programs.
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			t.Errorf("%s is not a Go file", name)
			continue
		}
		if strings.HasPrefix(name, "tag_") {
			// This file is used for tag processing in TestTags, below.
			continue
		}
		if name == "cgo.go" && !build.Default.CgoEnabled {
			t.Logf("cgo is not enabled for %s", name)
			continue
		}
		epointerCompileAndRun(t, t.TempDir(), epointer, name)
	}
}

// TestTags verifies that the -tags flag works as advertised.
func TestTags(t *testing.T) {
	epointer := epointerPath(t)
	dir := t.TempDir()
	var (
		protectedConstType = []byte("ProtectedConst")
		output             = filepath.Join(dir, "epointer_gen.go")
	)
	for _, file := range []string{"tag_main.go", "tag_tag.go"} {
		err := copy(filepath.Join(dir, file), filepath.Join("testdata", file))
		if err != nil {
			t.Fatal(err)
		}
	}
	// Run epointer in the directory that contains the package files.
	// We cannot run epointer in the current directory for the following reasons:
	// - Versions of Go earlier than Go 1.11, do not support absolute directories as a pattern.
	// - When the current directory is inside a go module, the path will not be considered
	//   a valid path to a package.
	err := runInDir(dir, epointer, ".")
	if err != nil {
		t.Fatal(err)
	}
	result, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(result, protectedConstType) {
		t.Fatal("tagged type appears in untagged run")
	}
	err = os.Remove(output)
	if err != nil {
		t.Fatal(err)
	}
	err = runInDir(dir, epointer, "-tags", "tag", ".")
	if err != nil {
		t.Fatal(err)
	}
	result, err = os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, protectedConstType) {
		t.Fatal("tagged type does not appear in tagged run")
	}
}

var exe struct {
	path string
	err  error
	once sync.Once
}

func epointerPath(t *testing.T) string {
	exe.once.Do(func() {
		exe.path, exe.err = os.Executable()
	})
	if exe.err != nil {
		t.Fatal(exe.err)
	}
	return exe.path
}

// epointerCompileAndRun runs epointer for the named file and compiles and
// runs the target binary in directory dir. That binary will panic if the Pointer methods are incorrect.
func epointerCompileAndRun(t *testing.T, dir, epointer, fileName string) {
	t.Helper()
	t.Logf("run: %s\n", fileName)
	source := filepath.Join(dir, path.Base(fileName))
	err := copy(source, filepath.Join("testdata", fileName))
	if err != nil {
		t.Fatalf("copying file to temporary directory: %s", err)
	}
	epointerSource := filepath.Join(dir, "epointer_gen.go")
	// Run epointer in temporary directory.
	err = run(epointer, "-output", epointerSource, source)
	if err != nil {
		t.Fatal(err)
	}
	// Run the binary in the temporary directory.
	err = run("go", "run", epointerSource, source)
	if err != nil {
		t.Fatal(err)
	}
}

// copy copies the from file to the to file.
func copy(to, from string) error {
	toFd, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFd.Close()
	fromFd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFd.Close()
	_, err = io.Copy(toFd, fromFd)
	return err
}

// run runs a single command and returns an error if it does not succeed.
// os/exec should have this function, to be honest.
func run(name string, arg ...string) error {
	return runInDir(".", name, arg...)
}

// runInDir runs a single command in directory dir and returns an error if
// it does not succeed.
func runInDir(dir, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GO111MODULE=auto")
	return cmd.Run()
}
