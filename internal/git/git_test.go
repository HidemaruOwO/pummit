package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/HidemaruOwO/pummit/internal/config"
)

// Tests in this file modify the working directory; do not run in parallel.

// Verify rune-safe truncation logic.
func TestTruncateWithEllipsis(t *testing.T) {
	t.Run("ascii", func(t *testing.T) {
		s := "abcde"
		got := truncateWithEllipsis(s, 3)
		want := "abc..."
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("multibyte", func(t *testing.T) {
		s := "あいうえお"
		got := truncateWithEllipsis(s, 3)
		want := "あいう..."
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("limit zero", func(t *testing.T) {
		s := "abc"
		got := truncateWithEllipsis(s, 0)
		if got != s {
			t.Fatalf("got %q, want %q", got, s)
		}
	})

	t.Run("exact limit", func(t *testing.T) {
		s := "abcd"
		got := truncateWithEllipsis(s, 4)
		if got != s {
			t.Fatalf("got %q, want %q", got, s)
		}
	})

	t.Run("emoji", func(t *testing.T) {
		// spec: truncates at rune boundaries, not grapheme clusters
		s := "✨✨"
		got := truncateWithEllipsis(s, 1)
		want := "✨..."
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})

	t.Run("negative limit", func(t *testing.T) {
		// spec: negative limit leaves string unchanged
		s := "abc"
		got := truncateWithEllipsis(s, -1)
		if got != s {
			t.Fatalf("got %q, want %q", got, s)
		}
	})
}

func TestCommitWithOfflineMode(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	tests := []struct {
		name         string
		emoji        string
		offline      bool
		aliasEnabled bool
		want         string
	}{
		{"online", "sparkles", false, false, "✨ test (file1.txt)"},
		{"offline", "sparkles", true, false, "✨ test (file1.txt)"},
		{"alias", "feat", true, true, "✨ test (file1.txt)"},
		{"alias fallback", "unknown", true, true, ":unknown: test (file1.txt)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setCleanGitEnv(t)

			dir := initTempRepo(t)
			if err := os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("a"), 0o644); err != nil {
				t.Fatalf("write file: %v", err)
			}
			runGit(t, dir, "add", "file1.txt")

			cwd, err := os.Getwd()
			if err != nil {
				t.Fatalf("getwd: %v", err)
			}
			if err := os.Chdir(dir); err != nil {
				t.Fatalf("chdir: %v", err)
			}
			t.Cleanup(func() { _ = os.Chdir(cwd) })

			orig := config.CurrentTOMLConfig
			t.Cleanup(func() { config.CurrentTOMLConfig = orig })
			config.CurrentTOMLConfig = config.GetDefaultTOMLConfig()
			config.CurrentTOMLConfig.Alias.Enabled = tt.aliasEnabled

			cm := CommitMessage{Emoji: tt.emoji, Message: "test"}
			if err := CommitWithOfflineMode(cm, tt.offline); err != nil {
				t.Fatalf("commit failed: %v", err)
			}

			msg := strings.TrimSpace(runGitOutput(t, dir, "log", "-1", "--pretty=%B"))
			if msg != tt.want {
				t.Fatalf("got %q, want %q", msg, tt.want)
			}
		})
	}
}

func TestGetChangedFilesList(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}

	cases := []struct {
		name  string
		files []string
	}{
		{"clean", nil},
		{"one file", []string{"a.txt"}},
		{"multiple files", []string{"a.txt", "b.txt"}},
		{"spaces and unicode", []string{"a b.txt", filepath.Join("サブ", "ほげ.txt")}},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			setCleanGitEnv(t)

			dir := initTempRepo(t)

			for _, f := range cs.files {
				path := filepath.Join(dir, f)
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				if err := os.WriteFile(path, []byte("a"), 0o644); err != nil {
					t.Fatalf("write file: %v", err)
				}
			}
			if len(cs.files) > 0 {
				args := append([]string{"add"}, cs.files...)
				runGit(t, dir, args...)
			}

			cwd, err := os.Getwd()
			if err != nil {
				t.Fatalf("getwd: %v", err)
			}
			if err := os.Chdir(dir); err != nil {
				t.Fatalf("chdir: %v", err)
			}
			t.Cleanup(func() { _ = os.Chdir(cwd) })

			normalizePath := func(s string) string {
				if unq, err := strconv.Unquote(s); err == nil {
					s = unq
				}
				return filepath.ToSlash(s)
			}

			got, err := GetChangedFilesList()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// spec: order-agnostic; compare as set
			if len(got) != len(cs.files) {
				t.Fatalf("got %v, want %v", got, cs.files)
			}
			wantSet := make(map[string]bool, len(cs.files))
			for _, f := range cs.files {
				wantSet[normalizePath(f)] = true
			}
			for _, f := range got {
				if !wantSet[normalizePath(f)] {
					t.Fatalf("unexpected file %q in %v", f, got)
				}
			}
		})
	}
}

func TestConvertToEmojiWithOfflineMode(t *testing.T) {
	t.Run("known", func(t *testing.T) {
		got := ConvertToEmojiWithOfflineMode("sparkles", false)
		if got != "✨" {
			t.Fatalf("got %q", got)
		}
	})

	t.Run("fallback offline", func(t *testing.T) {
		got := ConvertToEmojiWithOfflineMode("unknown", true)
		if got != ":unknown:" {
			t.Fatalf("got %q", got)
		}
	})

	t.Run("fallback online", func(t *testing.T) {
		got := ConvertToEmojiWithOfflineMode("unknown", false)
		if got != ":unknown:" {
			t.Fatalf("got %q", got)
		}
	})
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed: git %v: %v\n--- output ---\n%s", args, err, string(out))
	}
}

func runGitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed: git %v: %v\n--- output ---\n%s", args, err, string(out))
	}
	return string(out)
}

func initTempRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "init", "-q", "-b", "main")
	runGit(t, dir, "config", "user.name", "test")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "core.autocrlf", "false")
	return dir
}

// setCleanGitEnv isolates git from user-level config.
func setCleanGitEnv(t *testing.T) {
	t.Helper()
	t.Setenv("LC_ALL", "C")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_SYSTEM", "")
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
}
