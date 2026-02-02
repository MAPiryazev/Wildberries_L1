package cut

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	tests := []struct {
		name    string
		delim   string
		fields  []int
		wantErr bool
	}{
		{"valid", ",", []int{1, 2}, false},
		{"empty delimiter", "", []int{1}, true},
		{"empty fields", ",", []int{}, true},
		{"negative field", ",", []int{-1}, false},
		{"zero field", ",", []int{0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProcessor(tt.delim, tt.fields, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProcessor() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessLine(t *testing.T) {
	tests := []struct {
		name            string
		line            string
		delim           string
		fields          []int
		suppressNoDelim bool
		want            string
		wantErr         bool
	}{
		{
			name:    "single field",
			line:    "a,b,c",
			delim:   ",",
			fields:  []int{1},
			want:    "a",
			wantErr: false,
		},
		{
			name:    "multiple fields",
			line:    "a,b,c",
			delim:   ",",
			fields:  []int{1, 3},
			want:    "a,c",
			wantErr: false,
		},
		{
			name:    "field out of range",
			line:    "a,b",
			delim:   ",",
			fields:  []int{1, 5},
			want:    "a",
			wantErr: false,
		},
		{
			name:            "no delimiter suppress",
			line:            "abc",
			delim:           ",",
			fields:          []int{1},
			suppressNoDelim: true,
			want:            "",
			wantErr:         false,
		},
		{
			name:            "no delimiter no suppress",
			line:            "abc",
			delim:           ",",
			fields:          []int{1},
			suppressNoDelim: false,
			want:            "abc",
			wantErr:         false,
		},
		{
			name:    "custom delimiter",
			line:    "a:b:c",
			delim:   ":",
			fields:  []int{2},
			want:    "b",
			wantErr: false,
		},
		{
			name:    "colon delimiter",
			line:    "user:pass:domain",
			delim:   ":",
			fields:  []int{1, 3},
			want:    "user:domain",
			wantErr: false,
		},
		{
			name:    "space delimiter",
			line:    "one two three",
			delim:   " ",
			fields:  []int{2},
			want:    "two",
			wantErr: false,
		},
		{
			name:    "empty field",
			line:    "a,,c",
			delim:   ",",
			fields:  []int{2},
			want:    "",
			wantErr: false,
		},
		{
			name:    "field 0",
			line:    "a,b",
			delim:   ",",
			fields:  []int{0},
			wantErr: true,
		},
		{
			name:    "unsorted fields auto-sort",
			line:    "a,b,c,d",
			delim:   ",",
			fields:  []int{3, 1},
			want:    "a,c",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewProcessor(tt.delim, tt.fields, tt.suppressNoDelim)
			if err != nil {
				t.Fatalf("NewProcessor failed: %v", err)
			}

			got, err := p.ProcessLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessLine() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ProcessLine() got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessReader(t *testing.T) {
	input := "a,b,c\nd,e,f\ng,h,i\n"
	p, err := NewProcessor(",", []int{1, 3}, false)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	reader := strings.NewReader(input)
	var buf bytes.Buffer

	err = p.ProcessReader(reader, &buf)
	if err != nil {
		t.Errorf("ProcessReader() error = %v", err)
	}

	want := "a,c\nd,f\ng,i\n"
	got := buf.String()
	if got != want {
		t.Errorf("ProcessReader() got %q, want %q", got, want)
	}
}

func TestParseFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{"single field", "1", []int{1}, false},
		{"multiple fields", "1,2,3", []int{1, 2, 3}, false},
		{"range", "1-3", []int{1, 2, 3}, false},
		{"mixed", "1,3-5,7", []int{1, 3, 4, 5, 7}, false},
		{"duplicates", "1,2,1", []int{1, 2}, false},
		{"spaces", "1, 2, 3", []int{1, 2, 3}, false},
		{"empty", "", nil, true},
		{"invalid number", "a", nil, true},
		{"invalid range", "1-2-3", nil, true},
		{"reverse range", "5-1", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFields(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFields() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !equal(got, tt.want) {
				t.Errorf("ParseFields() got %v, want %v", got, tt.want)
			}
		})
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestProcessLineErrorHandling(t *testing.T) {
	p, _ := NewProcessor(",", []int{1}, false)

	_, err := p.ProcessLine("test")
	if !errors.Is(err, nil) && err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
