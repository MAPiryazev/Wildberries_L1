package unixcut_test

import (
	"reflect"
	"testing"

	"L2.13/pkg/unixcut"
)

func Test_WorkLines_Simple(t *testing.T) {
	lines := []string{"a,b,c", "1,2,3"}
	flags := unixcut.Flags{Delimiter: ","}
	expected := [][]string{{"a", "b", "c"}, {"1", "2", "3"}}

	out, err := unixcut.WorkLines(lines, flags)
	if err != nil {
		t.Errorf("%v", err)
	}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("ожидали: %v получили: %v", expected, out)
	}
}

func Test_WorkLines_SomeFlags_1(t *testing.T) {
	lines := []string{"a,b,c", "1,2,3"}
	flags := unixcut.Flags{Delimiter: ",", Fields: []int{1, 2}}
	expected := [][]string{{"b", "c"}, {"2", "3"}}

	out, err := unixcut.WorkLines(lines, flags)
	if err != nil {
		t.Errorf("%v", err)
	}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("ожидали: %v получили: %v", expected, out)
	}
}

func Test_WorkLines_SomeFlags_2(t *testing.T) {
	lines := []string{"a,b,c", "1,2,3"}
	flags := unixcut.Flags{Delimiter: ",", Fields: []int{1, 2}}
	expected := [][]string{{"b", "c"}, {"2", "3"}}

	out, err := unixcut.WorkLines(lines, flags)
	if err != nil {
		t.Errorf("%v", err)
	}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("ожидали: %v получили: %v", expected, out)
	}
}

func Test_WorkLines_SomeFlags_3(t *testing.T) {
	lines := []string{"a,b,c", "1,2,3"}
	flags := unixcut.Flags{Delimiter: ",", Fields: []int{}}
	expected := [][]string{{"a", "b", "c"}, {"1", "2", "3"}}

	out, err := unixcut.WorkLines(lines, flags)
	if err != nil {
		t.Errorf("%v", err)
	}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("ожидали: %v получили: %v", expected, out)
	}
}

func Test_WorkLines_SomeFlags_4(t *testing.T) {
	lines := []string{"a,b,c", "1,2,3", "noSep"}
	flags := unixcut.Flags{Delimiter: ",", Separated: true}
	expected := [][]string{{"a", "b", "c"}, {"1", "2", "3"}}

	out, err := unixcut.WorkLines(lines, flags)
	if err != nil {
		t.Errorf("%v", err)
	}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("ожидали: %v получили: %v", expected, out)
	}
}
