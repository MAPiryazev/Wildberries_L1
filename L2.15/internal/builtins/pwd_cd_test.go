package builtins

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCdAndPwd(t *testing.T) {
	start, _ := os.Getwd()
	defer func() { _ = os.Chdir(start) }()

	tmp := t.TempDir()
	Cd([]string{"cd", tmp})

	dir, _ := os.Getwd()
	if filepath.Clean(dir) != filepath.Clean(tmp) {
		t.Fatalf("got %q want %q", dir, tmp)
	}
}
