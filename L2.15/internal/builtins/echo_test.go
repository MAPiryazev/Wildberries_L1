package builtins

import (
	"io"
	"os"
	"testing"
)

func TestEcho(t *testing.T) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	Echo([]string{"echo", "hello", "world"})

	_ = w.Close()
	os.Stdout = old

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	_ = r.Close()

	if string(data) != "hello world\n" {
		t.Fatalf("got %q", string(data))
	}
}
