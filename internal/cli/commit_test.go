package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteCommitRequiresFlags(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Execute([]string{"commit", "test"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error for missing flags")
	}

	exitErr, ok := err.(*ExitError)
	if !ok {
		t.Fatalf("error type = %T, want *ExitError", err)
	}

	if exitErr.Code != userErrorCode {
		t.Fatalf("exit code = %d, want %d", exitErr.Code, userErrorCode)
	}
}

func TestExecuteCompatCommitWithoutMessageShowsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Execute([]string{"sparkles"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("stdout = %q, want help output", stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}
