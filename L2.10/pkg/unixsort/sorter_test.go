package unixsort_test

import (
	"reflect"
	"testing"

	"L2.10/pkg/unixsort"
)

func TestSortLines_SimpleLetters(t *testing.T) {
	input := []string{"d", "c", "b", "a"}
	output := unixsort.SortLines(input, unixsort.Flags{})
	waiting := []string{"a", "b", "c", "d"}

	if !reflect.DeepEqual(output, waiting) {
		t.Errorf("получили %v, хотели %v", output, waiting)
	}
}

func TestSortLines_Numbers(t *testing.T) {
	input := []string{"4184434", "432", "50", "1"}
	out := unixsort.SortLines(input, unixsort.Flags{Numeric: true})
	waiting := []string{"1", "50", "432", "4184434"}

	if !reflect.DeepEqual(out, waiting) {
		t.Errorf("получили %v ожидали %v", out, waiting)
	}
}

func TestSortLines_Suffix(t *testing.T) {
	input := []string{"1GB", "10mb", "100B"}
	out := unixsort.SortLines(input, unixsort.Flags{Human: true})
	waiting := []string{"100B", "10mb", "1GB"}

	if !reflect.DeepEqual(out, waiting) {
		t.Errorf("получили %v ожидали %v", out, waiting)
	}
}

func TestSortLines_ColumnAndNumber(t *testing.T) {
	input := []string{"a	100", "b	10", "c	1"}
	out := unixsort.SortLines(input, unixsort.Flags{Column: 2, Numeric: true})
	waiting := []string{"c	1", "b	10", "a	100"}

	if !reflect.DeepEqual(out, waiting) {
		t.Errorf("получили %v ожидали %v", out, waiting)
	}

}

func TestSortLines_Unique(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b"}
	out := unixsort.SortLines(input, unixsort.Flags{Unique: true})
	waiting := []string{"a", "b", "c"}

	if !reflect.DeepEqual(out, waiting) {
		t.Errorf("получили %v ожидали %v", out, waiting)
	}
}

func TestSortLines_Reverse(t *testing.T) {
	input := []string{"a", "b", "c", "d"}
	out := unixsort.SortLines(input, unixsort.Flags{Reverse: true})
	waiting := []string{"d", "c", "b", "a"}

	if !reflect.DeepEqual(out, waiting) {
		t.Errorf("получили %v ожидали %v", out, waiting)
	}
}
