package jsonparser

import (
	"testing"
)

func TestParser(t *testing.T) {
	testcases := []struct {
		input   string
		wantErr string
		output  Result
	}{{
		input:   "4085551212",
		wantErr: "",
		output:  Result{area: "408", part1: "555", part2: "1212"},
	}, {
		input:   "408-555-1212",
		wantErr: "",
		output:  Result{area: "408", part1: "555", part2: "1212"},
	}, {
		input:   "408-5551212",
		wantErr: "syntax error",
	}}
	for _, tc := range testcases {
		v, err := Parse([]byte(tc.input))
		var gotErr string
		if err != nil {
			gotErr = err.Error()
		}
		if gotErr != tc.wantErr {
			t.Errorf("%s err: %v, want %v", tc.input, gotErr, tc.wantErr)
		}
		if v != tc.output {
			t.Errorf("%s: %v, want %v", tc.input, v, tc.output)
		}
	}
}
