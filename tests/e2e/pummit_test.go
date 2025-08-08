package e2e

import (
  "context"
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  "testing"
  "time"

  "github.com/BurntSushi/toml"

  "github.com/HidemaruOwO/pummit/internal/config"
  "github.com/HidemaruOwO/pummit/internal/variable"
)

var binaryPath string

const testTimeout = 20 * time.Second

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

  defer func() {
    if err := os.Remove(binaryPath); err != nil {
      fmt.Fprintln(os.Stderr, err)
    }
  }()

  code := m.Run()
  os.Exit(code)
}

func initRepo(t *testing.T, dir string) {
  t.Helper()
  run := func(args ...string) {
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Dir = dir
    if out, err := cmd.CombinedOutput(); err != nil {
      t.Fatalf("%s: %v\n%s", strings.Join(args, " "), err, string(out))
    }
  }
  run("git", "init", "--initial-branch=main")
  run("git", "config", "user.email", "test@example.com")
  run("git", "config", "user.name", "Test")
}

func newTestRepo(t *testing.T) string {
  t.Helper()
  dir := t.TempDir()
  t.Setenv("HOME", dir)
  t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, ".config"))
  t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(dir, "nonexistent"))
  initRepo(t, dir)
  return dir
}

func TestCommit(t *testing.T) {
  dir := newTestRepo(t)

  file := filepath.Join(dir, "file.txt")
  if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
    t.Fatalf("write file: %v", err)
  }
  run := exec.Command("git", "add", "file.txt")
  run.Dir = dir
  if out, err := run.CombinedOutput(); err != nil {
    t.Fatalf("git add: %v\n%s", err, string(out))
  }

  ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
  defer cancel()
  cmd := exec.CommandContext(ctx, binaryPath, "sparkles", "test commit")
  cmd.Dir = dir
  if out, err := cmd.CombinedOutput(); err != nil {
    t.Fatalf("pummit commit: %v\n%s", err, string(out))
  }

  logCmd := exec.Command("git", "log", "-1", "--pretty=%B")
  logCmd.Dir = dir
  out, err := logCmd.CombinedOutput()
  if err != nil {
    t.Fatalf("git log: %v\n%s", err, string(out))
  }
  got := strings.TrimSpace(strings.ReplaceAll(string(out), "\r\n", "\n"))
  want1 := "✨ test commit (file.txt)"
  want2 := ":sparkles: test commit (file.txt)"
  if got != want1 && got != want2 {
    t.Fatalf("commit message = %q, want %q or %q", got, want1, want2)
  }
}

func TestAliasAdd(t *testing.T) {
  dir := newTestRepo(t)

  ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
  defer cancel()
  cmd := exec.CommandContext(
    ctx, binaryPath, "alias", "add", "rs", "rocket", "--emoji", "🚀",
  )
  cmd.Dir = dir
  if out, err := cmd.CombinedOutput(); err != nil {
    t.Fatalf("alias add: %v\n%s", err, string(out))
  }

  cfgPath := filepath.Join(
    os.Getenv("XDG_CONFIG_HOME"), "pummit", "config.toml",
  )
  data, err := os.ReadFile(cfgPath)
  if err != nil {
    t.Fatalf("read config: %v", err)
  }
  var cfg config.TOMLConfig
  if err := toml.Unmarshal(data, &cfg); err != nil {
    t.Fatalf("unmarshal config: %v", err)
  }

  found := false
search:
  for _, a := range cfg.Alias.Entries {
    for _, s := range a.Shortcuts {
      if s == "rs" && a.Name == "rocket" && a.Emoji == "🚀" {
        found = true
        break search
      }
    }
  }
  if !found {
    t.Fatalf("alias not added to config")
  }
}

func TestVersion(t *testing.T) {
  dir := newTestRepo(t)

  ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
  defer cancel()
  cmd := exec.CommandContext(ctx, binaryPath, "--version")
  cmd.Dir = dir
  out, err := cmd.CombinedOutput()
  if err != nil {
    t.Fatalf("version: %v\n%s", err, string(out))
  }
  got := strings.TrimSpace(strings.ReplaceAll(string(out), "\r\n", "\n"))
  want := fmt.Sprintf("pummit v%s", variable.VERSION)
  if got != want {
    t.Fatalf("version = %q, want %q", got, want)
  }
}
