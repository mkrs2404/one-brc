package main

import (
	"reflect"
	"testing"
)

func Test_parseBytes(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 int
	}{
		{
			name: "Positive temperature",
			args: args{
				line: []byte("Hamburg;12.9"),
			},
			want:  []byte("Hamburg"),
			want1: 129,
		},
		{
			name: "Negative temperature",
			args: args{
				line: []byte("Hamburg;-12.9"),
			},
			want:  []byte("Hamburg"),
			want1: -129,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseBytes(tt.args.line)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBytes() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseBytes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
