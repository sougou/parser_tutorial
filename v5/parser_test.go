package jsonparser

import (
	"testing"
)

func TestParser(t *testing.T) {
	testcases := []struct {
		input   string
		output  interface{}
		wantErr string
	}{{
		input:  `"abcd"`,
		output: "abcd",
	}, {
		input:  `"abc\t"`,
		output: "abc\t",
	}, {
		input:   `"abc\t`,
		wantErr: `syntax error`,
	}, {
		input:  `123`,
		output: float64(123),
	}, {
		input:  `123.1`,
		output: float64(123.1),
	}, {
		input:  `0.1`,
		output: float64(0.1),
	}, {
		input:  `123.1e-1`,
		output: float64(123.1e-1),
	}, {
		input:  `123.1e1`,
		output: float64(123.1e1),
	}, {
		input:   `.1`,
		wantErr: `syntax error`,
	}, {
		input:   `invalid`,
		wantErr: `syntax error`,
	}}
	for _, tc := range testcases {
		got, err := Parse([]byte(tc.input))
		var gotErr string
		if err != nil {
			gotErr = err.Error()
		}
		if gotErr != tc.wantErr {
			t.Errorf(`%s err: %v, want %v`, tc.input, gotErr, tc.wantErr)
		}
		if got != tc.output {
			t.Errorf(`%s: %v want %v`, tc.input, got, tc.output)
		}
	}
}
