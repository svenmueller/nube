package main

import (
	"io/ioutil"
	"testing"
)

func setup() {
}

func teardown() {
}

func TestAppWithoutApiKey(t *testing.T) {
	app := buildApp()
	app.Writer = ioutil.Discard

	// test with other global flags
	tests := []struct {
		args    []string
		wantErr error
	}{
		{
			args:    []string{"ctp-cli"},
			wantErr: nil,
		},
		{
			args:    []string{"ctp-cli", "--version"},
			wantErr: nil,
		},
		{
			args:    []string{"ctp-cli", "--help"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		err := app.Run(tt.args)
		if tt.wantErr == nil {
			if err != nil {
				t.Errorf("Expected nil error")
			}
		} else if err.Error() != tt.wantErr.Error() {
			t.Errorf("app.Run(%v) = %#v, want %#v", tt.args, err, tt.wantErr)
		}
	}
}

func TestGlobalFormatFlag(t *testing.T) {
	app := buildApp()
	app.Writer = ioutil.Discard

	args := []string{"ctp-cli", "-k", "key", "-f", "invalid"}
	err := app.Run(args)

	expected := `invalid output format: "invalid", available output options: json, yaml.`
	if err.Error() != expected {
		t.Errorf("app.Run(%v) = %#v, want %#v", err, err.Error(), expected)
	}
}
