package core

import (
	"io"
	"os"
	"testing"
)

func TestPipelineEchoUpper(t *testing.T) {
	shell := NewCore()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	shell.ExecuteLine("echo hello | tr a-z A-Z")

	_ = w.Close()
	os.Stdout = old

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	_ = r.Close()

	if string(data) != "HELLO\n" {
		t.Fatalf("got %q", string(data))
	}
}
