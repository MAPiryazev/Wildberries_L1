package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "simple with repeats",
			input:   "lkjgklerg4",
			want:    "lkjgklergggg",
			wantErr: false,
		},
		{
			name:    "no digits",
			input:   "abcd",
			want:    "abcd",
			wantErr: false,
		},
		{
			name:    "only digits invalid",
			input:   "45",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "escaped digits separately",
			input:   "qwe\\4\\5",
			want:    "qwe45",
			wantErr: false,
		},
		{
			name:    "escaped digit then repeat",
			input:   "qwe\\45",
			want:    "qwe44444",
			wantErr: false,
		},
		{
			name:    "ending with slash invalid",
			input:   "abc\\",
			want:    "",
			wantErr: true,
		},
		{
			name:    "digit zero removes symbol",
			input:   "a0bc",
			want:    "bc",
			wantErr: false,
		},
		{
			name:    "digit one keeps symbol",
			input:   "a1b",
			want:    "ab",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpackSrting(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected error: %v, wantErr=%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("input=%q got=%q, want=%q", tt.input, got, tt.want)
			}
		})
	}

}
