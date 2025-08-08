package doctor

// NOTE: These tests mutate package-level globals in config.
// Do not use t.Parallel() in this file.

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/HidemaruOwO/pummit/internal/config"
)

// To accommodate variations in Git output (Apple Git / Windows derivatives), allow the numeric core plus any non-whitespace suffix.
var reGitVerAny = regexp.MustCompile(`\b(\d+\.\d+(?:\.\d+)?)(?:\S*)\b`)

// setCleanConfigEnvs removes host settings that could leak into tests.
func setCleanConfigEnvs(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
	if runtime.GOOS == "windows" {
		// Windows Git resolves HOME/USERPROFILE differently depending on environment.
		// Keep them aligned to avoid accidental leakage from the host.
		t.Setenv("USERPROFILE", home)
	}
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("LC_ALL", "C")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_SYSTEM", "")
	if runtime.GOOS == "windows" {
		t.Setenv("GIT_CONFIG_GLOBAL", "NUL")
	} else {
		t.Setenv("GIT_CONFIG_GLOBAL", "/dev/null")
	}
}

// writeGitConfig provides minimal identity so GitConfigChecker passes.
func writeGitConfig(t *testing.T, home string) {
	t.Helper()
	path := filepath.Join(home, ".gitconfig")
	data := "[user]\n\tname = Test\n\temail = test@example.com\n"
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatalf("write gitconfig: %v", err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", path)
}

// restoreConfigState avoids cross-test pollution of package globals.
func restoreConfigState(t *testing.T) {
	prev := config.CurrentTOMLConfig
	prevPath := config.TOMLConfigPath
	t.Cleanup(func() {
		config.CurrentTOMLConfig = prev
		config.TOMLConfigPath = prevPath
	})
}

// initRealRepo creates a minimal git repository and switches cwd.
func initRealRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	requireGit(t)
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v, %s", err, out)
	}
	wd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
	return dir
}

// Common prerequisite: check for Git (skip if not found).
func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
}

// Simple wrapper: a safety net to prevent runaway commands (protection against long hangs).
func runCmd(ctx context.Context, name string, args ...string) ([]byte, error) {
	c := exec.CommandContext(ctx, name, args...)
	return c.CombinedOutput()
}

func TestRunAllChecksHealthyEnv(t *testing.T) {
	requireGit(t)
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	writeGitConfig(t, home)
	restoreConfigState(t)

	// Run under a real, minimal git repository to satisfy repo checks.
	_ = initRealRepo(t)

	if err := config.LoadTOMLConfig(); err != nil {
		t.Fatalf("load config: %v", err)
	}
	results := RunAllChecks()
	if results.HasErrors() {
		t.Fatalf("want no errors, got %+v", results)
	}
}

func TestGetSystemInfoGitMissing(t *testing.T) {
	tmp := t.TempDir()
	setCleanConfigEnvs(t, tmp)

	// Remove git from PATH to simulate absence.
	t.Setenv("PATH", tmp)

	// Treat failure to detect the version string (numeric part + suffix) as indicating that Git is unavailable. Does not rely on the specifics of the message.
	info := GetSystemInfo()
	if reGitVerAny.MatchString(info.GitVersion) {
		t.Fatalf("expected no git version when git is missing, got %q", info.GitVersion)
	}
}

func TestRunAllChecksMalformedConfig(t *testing.T) {
	requireGit(t)
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	writeGitConfig(t, home)

	// Ensure repo-related checks are satisfied so only config failure surfaces.
	_ = initRealRepo(t)

	cfgDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatalf("mkdir cfg: %v", err)
	}
	bad := filepath.Join(cfgDir, "config.toml")
	if err := os.WriteFile(bad, []byte("invalid = ["), 0644); err != nil {
		t.Fatalf("write bad config: %v", err)
	}
	restoreConfigState(t)

	results := RunAllChecks()
	var cfgErr bool
	for _, r := range results {
		if r.Name == "Configuration Files" && r.Status == "ERROR" {
			cfgErr = true
			break
		}
	}
	if !cfgErr {
		t.Fatalf("expected config error, got %+v", results)
	}
}

func TestGetSystemInfoGitVersionFormat(t *testing.T) {
	requireGit(t)
	home := t.TempDir()
	setCleanConfigEnvs(t, home)

	info := GetSystemInfo()

	// Check that the Git version string contains a numeric part followed by any non-whitespace suffix.
	m := reGitVerAny.FindStringSubmatch(info.GitVersion)
	if m == nil {
		t.Fatalf("unexpected git version: %q", info.GitVersion)
	}
}
