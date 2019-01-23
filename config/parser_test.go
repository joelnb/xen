package config

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
		input:   "name = 'cloud'",
		wantErr: "",
	}, {
		input:   "name = \"cloud'",
		wantErr: "syntax error",
	}, {
		input:   `on_poweroff = 'destroy'; on_crash    = "preserve"`,
		wantErr: "",
	}, {
		input:   `vif = []`,
		wantErr: "",
	}, {
		input:   `vif = ["some-string",'another string']`,
		wantErr: "",
	}, {
		input:   `vif = ["some-string",]`,
		wantErr: "",
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
