package jsonparser

import (
	"testing"
)

func TestParser(t *testing.T) {
	testcases := []struct {
		input   string
		wantErr string
	}{{
		input:   "",
		wantErr: "",
	}, {
		input:   "1",
		wantErr: "syntax error",
	}}
	for _, tc := range testcases {
		_, err := Parse([]byte(tc.input))
		var gotErr string
		if err != nil {
			gotErr = err.Error()
		}
		if gotErr != tc.wantErr {
			t.Errorf("%s err: %v, want %v", tc.input, gotErr, tc.wantErr)
		}
	}
}
