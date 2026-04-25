package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/HidemaruOwO/pummit/legacy/config"
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
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	build := exec.Command("go", "build", "-o", binaryPath)
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		os.Exit(1)
	}

	code := func() int {
		defer func() {
			if err := os.Remove(binaryPath); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()
		return m.Run()
	}()

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
	t.Setenv("USERPROFILE", dir)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, ".config"))
	t.Setenv("APPDATA", filepath.Join(dir, "AppData", "Roaming"))
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

func TestCommitCommand(t *testing.T) {
	dir := newTestRepo(t)

	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	run := exec.Command("git", "add", "file.txt")
	run.Dir = dir
	if out, err := run.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, string(out))
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, binaryPath, "commit", "--emoji", "sparkles", "test commit")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pummit commit command: %v\n%s", err, string(out))
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

func TestAutoEmojiCommit(t *testing.T) {
	dir := newTestRepo(t)
	checkout := exec.Command("git", "checkout", "-b", "feature/test")
	checkout.Dir = dir
	if out, err := checkout.CombinedOutput(); err != nil {
		t.Fatalf("git checkout: %v\n%s", err, string(out))
	}

	file := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(file, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	run := exec.Command("git", "add", "file.txt")
	run.Dir = dir
	if out, err := run.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, string(out))
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, binaryPath, "commit", "--auto-emoji", "test commit")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pummit auto commit: %v\n%s", err, string(out))
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

	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("get config dir: %v", err)
	}

	cfgPath := filepath.Join(configDir, "config.toml")
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

func TestAliasDeleteAndReset(t *testing.T) {
	dir := newTestRepo(t)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	add := exec.CommandContext(ctx, binaryPath, "alias", "add", "rs", "rocket", "--emoji", "🚀")
	add.Dir = dir
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("alias add: %v\n%s", err, string(out))
	}

	deleteCmd := exec.CommandContext(ctx, binaryPath, "alias", "delete", "rs", "--force")
	deleteCmd.Dir = dir
	if out, err := deleteCmd.CombinedOutput(); err != nil {
		t.Fatalf("alias delete: %v\n%s", err, string(out))
	}

	list := exec.CommandContext(ctx, binaryPath, "alias", "list")
	list.Dir = dir
	out, err := list.CombinedOutput()
	if err != nil {
		t.Fatalf("alias list: %v\n%s", err, string(out))
	}
	if strings.Contains(string(out), "rs") {
		t.Fatalf("alias list still contains deleted shortcut: %s", string(out))
	}

	reset := exec.CommandContext(ctx, binaryPath, "alias", "reset", "--force")
	reset.Dir = dir
	if out, err := reset.CombinedOutput(); err != nil {
		t.Fatalf("alias reset: %v\n%s", err, string(out))
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

func TestVersionCommand(t *testing.T) {
	dir := newTestRepo(t)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, binaryPath, "version")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version command: %v\n%s", err, string(out))
	}
	got := strings.TrimSpace(strings.ReplaceAll(string(out), "\r\n", "\n"))
	want := fmt.Sprintf("pummit v%s", variable.VERSION)
	if got != want {
		t.Fatalf("version = %q, want %q", got, want)
	}
}

func TestConfigValidate(t *testing.T) {
	dir := newTestRepo(t)

	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("get config dir: %v", err)
	}

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte("[base]\nemoji=false\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, binaryPath, "config", "validate")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("config validate: %v\n%s", err, string(out))
	}

	got := strings.TrimSpace(strings.ReplaceAll(string(out), "\r\n", "\n"))
	if got != "Configuration is valid" {
		t.Fatalf("config validate output = %q, want success message", got)
	}
}

func TestConfigSetGetReset(t *testing.T) {
	dir := newTestRepo(t)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	setCmd := exec.CommandContext(ctx, binaryPath, "config", "set", "base.emoji", "false")
	setCmd.Dir = dir
	if out, err := setCmd.CombinedOutput(); err != nil {
		t.Fatalf("config set: %v\n%s", err, string(out))
	}

	getCmd := exec.CommandContext(ctx, binaryPath, "config", "get", "base.emoji")
	getCmd.Dir = dir
	out, err := getCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("config get: %v\n%s", err, string(out))
	}
	if strings.TrimSpace(string(out)) != "false" {
		t.Fatalf("config get output = %q, want false", strings.TrimSpace(string(out)))
	}

	resetCmd := exec.CommandContext(ctx, binaryPath, "config", "reset", "--force")
	resetCmd.Dir = dir
	if out, err := resetCmd.CombinedOutput(); err != nil {
		t.Fatalf("config reset: %v\n%s", err, string(out))
	}

	getAfter := exec.CommandContext(ctx, binaryPath, "config", "get", "base.emoji")
	getAfter.Dir = dir
	out, err = getAfter.CombinedOutput()
	if err != nil {
		t.Fatalf("config get after reset: %v\n%s", err, string(out))
	}
	if strings.TrimSpace(string(out)) != "true" {
		t.Fatalf("config get after reset = %q, want true", strings.TrimSpace(string(out)))
	}
}

func TestDoctor(t *testing.T) {
	dir := newTestRepo(t)

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, binaryPath, "doctor")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("doctor: %v\n%s", err, string(out))
	}

	text := string(out)
	if !strings.Contains(text, "Config") || !strings.Contains(text, "Git") || !strings.Contains(text, "Network") {
		t.Fatalf("doctor output missing expected sections: %s", text)
	}
}

func TestMigrateStatusAndRollback(t *testing.T) {
	dir := newTestRepo(t)

	configDir, err := config.GetConfigDir()
	if err != nil {
		t.Fatalf("get config dir: %v", err)
	}
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config dir: %v", err)
	}
	jsonPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(jsonPath, []byte(variable.DEFAULT_CONFIG), 0o644); err != nil {
		t.Fatalf("write config.json: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	migrateCmd := exec.CommandContext(ctx, binaryPath, "migrate")
	migrateCmd.Dir = dir
	if out, err := migrateCmd.CombinedOutput(); err != nil {
		t.Fatalf("migrate: %v\n%s", err, string(out))
	}

	statusCmd := exec.CommandContext(ctx, binaryPath, "migrate", "status")
	statusCmd.Dir = dir
	out, err := statusCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("migrate status: %v\n%s", err, string(out))
	}
	if !strings.Contains(string(out), "status: toml_only") {
		t.Fatalf("unexpected migrate status output: %s", string(out))
	}

	backupPath := filepath.Join(configDir, "config.json.bak")
	if _, err := os.Stat(backupPath); err != nil {
		t.Fatalf("backup file missing: %v", err)
	}

	rollbackCmd := exec.CommandContext(ctx, binaryPath, "migrate", "rollback", "--confirm")
	rollbackCmd.Dir = dir
	if out, err := rollbackCmd.CombinedOutput(); err != nil {
		t.Fatalf("rollback: %v\n%s", err, string(out))
	}

	if _, err := os.Stat(jsonPath); err != nil {
		t.Fatalf("config.json not restored: %v", err)
	}
	if _, err := os.Stat(filepath.Join(configDir, "config.toml")); !os.IsNotExist(err) {
		t.Fatalf("config.toml should be removed after rollback")
	}
}
