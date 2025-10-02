package unixgrep_test

import (
	"reflect"
	"testing"

	"L2.12/pkg/unixgrep"
)

func TestHandleGrep_FixedString(t *testing.T) {
	lines := []string{
		"hello world",
		"HELLO WORLD",
		"goodbye",
	}
	flags := unixgrep.Flags{
		FixedString: true,
	}

	got := unixgrep.HandleGrep(lines, "hello", flags)
	want := []string{"hello world"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ожидали %v, получили %v", want, got)
	}
}

func TestHandleGrep_IgnoreCase(t *testing.T) {
	lines := []string{
		"hello world",
		"HELLO WORLD",
		"goodbye",
	}
	flags := unixgrep.Flags{
		IgnoreCase: true,
	}

	got := unixgrep.HandleGrep(lines, "hello", flags)
	want := []string{"hello world", "HELLO WORLD"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ожидали %v, получили %v", want, got)
	}
}

func TestHandleGrep_CountOnly(t *testing.T) {
	lines := []string{
		"aaa",
		"bbb",
		"aaa bbb",
		"ccc",
	}
	flags := unixgrep.Flags{
		FixedString: true,
		CountOnly:   true,
	}

	got := unixgrep.HandleGrep(lines, "aaa", flags)
	want := []string{"2"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ожидали %v, получили %v", want, got)
	}
}
