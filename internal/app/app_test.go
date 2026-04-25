package app

import (
	"errors"
	"testing"
)

func TestShouldUseNewCLI(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "no args", args: nil, want: true},
		{name: "version flag", args: []string{"--version"}, want: true},
		{name: "version command", args: []string{"version"}, want: true},
		{name: "help command", args: []string{"help"}, want: true},
		{name: "help flag", args: []string{"--help"}, want: true},
		{name: "config validate", args: []string{"config", "validate"}, want: true},
		{name: "config list new", args: []string{"config", "list"}, want: true},
		{name: "new commit command", args: []string{"commit", "--emoji", "sparkles", "test"}, want: true},
		{name: "compat commit args", args: []string{"sparkles", "test"}, want: true},
		{name: "alias subcommand", args: []string{"alias", "list"}, want: true},
		{name: "doctor command", args: []string{"doctor"}, want: true},
		{name: "migrate command", args: []string{"migrate", "status"}, want: true},
		{name: "new mcp", args: []string{"mcp"}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldUseNewCLI(tt.args)
			if got != tt.want {
				t.Fatalf("shouldUseNewCLI(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "nil", err: nil, want: ExitCodeSuccess},
		{name: "plain error", err: errors.New("boom"), want: ExitCodeError},
		{name: "coded error", err: NewExitError(7, errors.New("boom")), want: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExitCode(tt.err)
			if got != tt.want {
				t.Fatalf("ExitCode(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}
