package logger

import (
	"io"
	"os"
	"testing"

	"github.com/fatih/color"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	oldColorOutput := color.Output
	oldNoColor := color.NoColor

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w
	color.Output = w
	color.NoColor = true

	defer func() {
		os.Stdout = oldStdout
		color.Output = oldColorOutput
		color.NoColor = oldNoColor
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("failed to close reader: %v", err)
	}

	return string(output)
}

func TestNewDebugFlag(t *testing.T) {
	t.Run("debug enabled", func(t *testing.T) {
		t.Setenv("DEBUG", "1")

		logger := New()

		if !logger.debug {
			t.Fatalf("expected debug mode when DEBUG=1")
		}
	})

	t.Run("debug disabled", func(t *testing.T) {
		t.Setenv("DEBUG", "")

		logger := New()

		if logger.debug {
			t.Fatalf("expected debug mode to be disabled when DEBUG is unset")
		}
	})
}

func TestInfo(t *testing.T) {
	t.Setenv("DEBUG", "")
	logger := New()

	got := captureOutput(t, func() {
		logger.Info("info message")
	})

	want := "info message\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestInfof(t *testing.T) {
	t.Setenv("DEBUG", "")
	logger := New()

	got := captureOutput(t, func() {
		logger.Infof("info %s", "message")
	})

	want := "info message\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestError(t *testing.T) {
	t.Setenv("DEBUG", "")
	logger := New()

	got := captureOutput(t, func() {
		logger.Error("error message")
	})

	want := "error message\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestErrorf(t *testing.T) {
	t.Setenv("DEBUG", "")
	logger := New()

	got := captureOutput(t, func() {
		logger.Errorf("error %d", 1)
	})

	want := "error 1\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestSuccess(t *testing.T) {
	t.Setenv("DEBUG", "")
	logger := New()

	got := captureOutput(t, func() {
		logger.Success("success message")
	})

	want := "success message\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestSuccessf(t *testing.T) {
	t.Setenv("DEBUG", "")
	logger := New()

	got := captureOutput(t, func() {
		logger.Successf("success %s", "message")
	})

	want := "success message\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestDebug(t *testing.T) {
	t.Run("debug enabled outputs", func(t *testing.T) {
		t.Setenv("DEBUG", "1")
		logger := New()

		got := captureOutput(t, func() {
			logger.Debug("debug message")
		})

		want := "[DEBUG] debug message\n"
		if got != want {
			t.Fatalf("unexpected output: got %q, want %q", got, want)
		}
	})

	t.Run("debug disabled no output", func(t *testing.T) {
		t.Setenv("DEBUG", "")
		logger := New()

		got := captureOutput(t, func() {
			logger.Debug("debug message")
		})

		if got != "" {
			t.Fatalf("expected no output when debug disabled, got %q", got)
		}
	})
}

func TestDebugf(t *testing.T) {
	t.Run("debug enabled outputs", func(t *testing.T) {
		t.Setenv("DEBUG", "1")
		logger := New()

		got := captureOutput(t, func() {
			logger.Debugf("value=%d", 42)
		})

		want := "[DEBUG] value=42\n"
		if got != want {
			t.Fatalf("unexpected output: got %q, want %q", got, want)
		}
	})

	t.Run("debug disabled no output", func(t *testing.T) {
		t.Setenv("DEBUG", "")
		logger := New()

		got := captureOutput(t, func() {
			logger.Debugf("value=%d", 42)
		})

		if got != "" {
			t.Fatalf("expected no output when debug disabled, got %q", got)
		}
	})
}
