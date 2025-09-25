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
			name:    "повторения",
			input:   "lkjgklerg4",
			want:    "lkjgklergggg",
			wantErr: false,
		},
		{
			name:    "нет цифр",
			input:   "abcd",
			want:    "abcd",
			wantErr: false,
		},
		{
			name:    "только цифры",
			input:   "45",
			want:    "",
			wantErr: true,
		},
		{
			name:    "пустая строка",
			input:   "",
			want:    "",
			wantErr: false,
		},
		{
			name:    "разделенные цифры",
			input:   "qwe\\4\\5",
			want:    "qwe45",
			wantErr: false,
		},
		{
			name:    "цифры с повторением",
			input:   "qwe\\45",
			want:    "qwe44444",
			wantErr: false,
		},
		{
			name:    "еще кейс",
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
