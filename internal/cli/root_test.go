package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/HidemaruOwO/pummit/internal/variable"
)

func TestExecuteVersionFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Execute([]string{"--version"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	want := fmt.Sprintf("pummit v%s\n", variable.VERSION)
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestExecuteVersionCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Execute([]string{"version"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	want := fmt.Sprintf("pummit v%s\n", variable.VERSION)
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestExecuteHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Execute(nil, &stdout, &stderr)
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
