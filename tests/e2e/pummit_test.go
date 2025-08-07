package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/internal/variable"
)

var binaryPath string

func TestMain(m *testing.M) {
	root, err := filepath.Abs("../..")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	binaryPath = filepath.Join(root, "pummit-e2e")
	build := exec.Command("go", "build", "-o", binaryPath)
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		os.Exit(1)
	}

	code := m.Run()
	os.Remove(binaryPath)
	os.Exit(code)
}

func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", strings.Join(append([]string{name}, args...), " "), err, string(out))
	}
	return string(out)
}

func initRepo(t *testing.T, dir string) {
	t.Helper()
	run(t, dir, "git", "init", "--initial-branch=main")
	run(t, dir, "git", "config", "user.email", "test@example.com")
	run(t, dir, "git", "config", "user.name", "Test")
}

func TestCommit(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, ".config"))
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(dir, "nonexistent"))
	initRepo(t, dir)

	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	run(t, dir, "git", "add", "file.txt")

	run(t, dir, binaryPath, "sparkles", "test commit")

	out := run(t, dir, "git", "log", "-1", "--pretty=%B")
	got := strings.TrimSpace(strings.ReplaceAll(out, "\r\n", "\n"))
	want1 := "✨ test commit (file.txt)"
	want2 := ":sparkles: test commit (file.txt)"
	if got != want1 && got != want2 {
		t.Fatalf("commit message = %q, want %q or %q", got, want1, want2)
	}
}

func TestAliasAdd(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, ".config"))
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(dir, "nonexistent"))
	initRepo(t, dir)

	// NOTE: The command syntax is expected to change to "alias add" in the future.
	run(t, dir, binaryPath, "alias:add", "rs", "rocket", "--emoji", "🚀")

	cfgPath := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "pummit", "config.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	var cfg config.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}

	found := false
	for _, a := range cfg.Aliases {
		names := strings.Split(a[0], ",")
		for _, n := range names {
			if n == "rs" && len(a) == 3 && a[1] == "rocket" && a[2] == "🚀" {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("alias not added to config")
	}
}

func TestVersion(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, ".config"))
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(dir, "nonexistent"))
	initRepo(t, dir)

	out := run(t, dir, binaryPath, "--version")
	got := strings.TrimSpace(strings.ReplaceAll(out, "\r\n", "\n"))
	want := fmt.Sprintf("pummit v%s", variable.VERSION)
	if got != want {
		t.Fatalf("version = %q, want %q", got, want)
	}
}
