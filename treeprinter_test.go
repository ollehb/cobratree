package cobratree

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/spf13/cobra"
)

func TestOneLevel(t *testing.T) {
	goldenContents, err := ioutil.ReadFile("_golden/1level")
	if err != nil {
		t.Fatalf("could not read golden file: %v", err)
	}

	root := &cobra.Command{Use: "root"}
	child1 := &cobra.Command{Use: "child1"}
	child2 := &cobra.Command{Use: "child2"}
	child3 := &cobra.Command{Use: "child3"}

	root.AddCommand(child1, child2, child3)
	buffer := new(bytes.Buffer)
	if err := WriteTree(buffer, root); err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(goldenContents, buffer.Bytes()) != 0 {
		t.Fatal("output does not match golden file")
	}
}

func TestComplexHierarchy(t *testing.T) {
	goldenContents, err := ioutil.ReadFile("_golden/complex")
	if err != nil {
		t.Fatalf("could not read golden file: %v", err)
	}

	root := &cobra.Command{Use: "root"}
	child1 := &cobra.Command{Use: "child1"}
	child2 := &cobra.Command{Use: "child2"}
	child3 := &cobra.Command{Use: "child3"}
	child4 := &cobra.Command{Use: "child4"}

	subchild := &cobra.Command{Use: "subchild"}
	subsibling := &cobra.Command{Use: "subsibling"}

	child1.AddCommand(subchild)
	child2.AddCommand(subsibling, subchild)

	subchild2 := &cobra.Command{Use: "subchild2"}
	subchild3 := &cobra.Command{Use: "subchild3"}
	subchild4 := &cobra.Command{Use: "subchild4"}
	subchild5 := &cobra.Command{Use: "subchild5"}

	child3.AddCommand(subchild2)
	subchild2.AddCommand(subchild3)
	subchild3.AddCommand(subchild4)
	subchild3.AddCommand(subsibling)
	subchild4.AddCommand(subchild5)

	root.AddCommand(child1, child2, child3, child4)
	buffer := new(bytes.Buffer)
	if err := WriteTree(buffer, root); err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(goldenContents, buffer.Bytes()) != 0 {
		t.Fatal("output does not match golden file")
	}
}

func TestRootOnly(t *testing.T) {
	commandUse := "root"
	root := &cobra.Command{Use: commandUse}
	expected := commandUse + "\n"

	out := new(bytes.Buffer)
	if err := WriteTree(out, root); err != nil {
		t.Fatal(err)
	}

	outString := out.String()
	if expected != outString {
		t.Fatalf("%s differs from %s", expected, outString)
	}
}

func TestWriterError(t *testing.T) {
	command := &cobra.Command{Use: "f"}
	if err := WriteTree(errorWriter{}, command); err == nil {
		t.Fatal("expected error")
	}
}

type errorWriter struct {
}

func (errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("always errors")
}
