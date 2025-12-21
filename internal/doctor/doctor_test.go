package doctor

// NOTE: These tests mutate package-level globals in config.
// Do not use t.Parallel() in this file.

// Parallel execution is problematic because package-level global mutation can cause
// race conditions and unpredictable test failures if tests run concurrently.
// All tests in this file must run serially to avoid cross-test interference.

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
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
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Logf("failed to restore working directory: %v", err)
		}
	})
	return dir
}

// Common prerequisite: check for Git (skip if not found).
func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
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

func TestDiagnosticResultsHasErrors(t *testing.T) {
	tests := []struct {
		name    string
		results DiagnosticResults
		want    bool
	}{
		{"no errors", DiagnosticResults{{Status: "OK"}, {Status: "WARNING"}}, false},
		{"has error", DiagnosticResults{{Status: "OK"}, {Status: "ERROR"}}, true},
		{"empty", DiagnosticResults{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.results.HasErrors(); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiagnosticResultsHasWarnings(t *testing.T) {
	tests := []struct {
		name    string
		results DiagnosticResults
		want    bool
	}{
		{"no warnings", DiagnosticResults{{Status: "OK"}, {Status: "ERROR"}}, false},
		{"has warning", DiagnosticResults{{Status: "OK"}, {Status: "WARNING"}}, true},
		{"empty", DiagnosticResults{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.results.HasWarnings(); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitConfigCheckerName(t *testing.T) {
	c := &GitConfigChecker{}
	if got := c.Name(); got != "Git Configuration" {
		t.Fatalf("got %q, want %q", got, "Git Configuration")
	}
}

func TestGitRepositoryCheckerName(t *testing.T) {
	c := &GitRepositoryChecker{}
	if got := c.Name(); got != "Git Repository Status" {
		t.Fatalf("got %q, want %q", got, "Git Repository Status")
	}
}

func TestGitRepositoryCheckerNotARepo(t *testing.T) {
	tmp := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	c := &GitRepositoryChecker{}
	result := c.Check()
	if result.Status != "WARNING" {
		t.Fatalf("expected WARNING, got %s", result.Status)
	}
}

func TestConfigFileCheckerName(t *testing.T) {
	c := &ConfigFileChecker{}
	if got := c.Name(); got != "Configuration Files" {
		t.Fatalf("got %q, want %q", got, "Configuration Files")
	}
}

func TestNetworkCheckerName(t *testing.T) {
	c := &NetworkChecker{}
	if got := c.Name(); got != "Network Connectivity" {
		t.Fatalf("got %q, want %q", got, "Network Connectivity")
	}
}

func TestFilePermissionCheckerName(t *testing.T) {
	c := &FilePermissionChecker{}
	if got := c.Name(); got != "File Permissions" {
		t.Fatalf("got %q, want %q", got, "File Permissions")
	}
}

func TestGetGoVersion(t *testing.T) {
	got := getGoVersion()
	if strings.HasPrefix(got, "go") {
		t.Fatalf("expected version without 'go' prefix, got %q", got)
	}
	matched, err := regexp.MatchString(`^\d+\.\d+`, got)
	if err != nil {
		t.Fatalf("match version: %v", err)
	}
	if !matched {
		t.Fatalf("unexpected version format: %q", got)
	}
}

func TestSystemInfoFormatSystemInfo(t *testing.T) {
	info := SystemInfo{
		OS:            "linux",
		Architecture:  "amd64",
		GoVersion:     "1.21.0",
		PummitVersion: "3.0.0",
		GitVersion:    "2.40.0",
	}

	got := info.FormatSystemInfo()

	if !strings.Contains(got, "linux") {
		t.Error("missing OS")
	}
	if !strings.Contains(got, "amd64") {
		t.Error("missing architecture")
	}
	if !strings.Contains(got, "1.21.0") {
		t.Error("missing Go version")
	}
	if !strings.Contains(got, "3.0.0") {
		t.Error("missing Pummit version")
	}
	if !strings.Contains(got, "2.40.0") {
		t.Error("missing Git version")
	}
}

func TestConfigFileCheckerNoConfigDir(t *testing.T) {
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	restoreConfigState(t)

	c := &ConfigFileChecker{}
	result := c.Check()

	if result.Status != "WARNING" {
		t.Fatalf("expected WARNING, got %s", result.Status)
	}
}

func TestConfigFileCheckerNoConfigFiles(t *testing.T) {
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	restoreConfigState(t)

	cfgDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatalf("mkdir cfg: %v", err)
	}

	c := &ConfigFileChecker{}
	result := c.Check()

	if result.Status != "WARNING" {
		t.Fatalf("expected WARNING, got %s", result.Status)
	}
}

func TestConfigFileCheckerValidJSONConfig(t *testing.T) {
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	restoreConfigState(t)

	cfgDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatalf("mkdir cfg: %v", err)
	}

	jsonPath := filepath.Join(cfgDir, "config.json")
	if err := os.WriteFile(jsonPath, []byte(`{"writeEmoji": true}`), 0644); err != nil {
		t.Fatalf("write json: %v", err)
	}

	c := &ConfigFileChecker{}
	result := c.Check()

	if result.Status != "OK" {
		t.Fatalf("expected OK, got %s: %s", result.Status, result.Message)
	}
}

func TestConfigFileCheckerInvalidJSONConfig(t *testing.T) {
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	restoreConfigState(t)

	cfgDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatalf("mkdir cfg: %v", err)
	}

	jsonPath := filepath.Join(cfgDir, "config.json")
	if err := os.WriteFile(jsonPath, []byte(`{invalid`), 0644); err != nil {
		t.Fatalf("write json: %v", err)
	}

	c := &ConfigFileChecker{}
	result := c.Check()

	if result.Status != "ERROR" {
		t.Fatalf("expected ERROR, got %s", result.Status)
	}
}

func TestGitConfigCheckerCheckMissing(t *testing.T) {
	requireGit(t)
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	restoreConfigState(t)

	emptyCfg := filepath.Join(home, ".gitconfig")
	if err := os.WriteFile(emptyCfg, nil, 0600); err != nil {
		t.Fatalf("write empty gitconfig: %v", err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", emptyCfg)

	c := &GitConfigChecker{}
	result := c.Check()

	if result.Status != "ERROR" {
		t.Fatalf("expected ERROR, got %s", result.Status)
	}
	if len(result.Suggestions) == 0 {
		t.Fatalf("expected suggestions for missing config")
	}
}

func TestGitConfigCheckerCheckOK(t *testing.T) {
	requireGit(t)
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	writeGitConfig(t, home)
	restoreConfigState(t)

	c := &GitConfigChecker{}
	result := c.Check()

	if result.Status != "OK" {
		t.Fatalf("expected OK, got %s", result.Status)
	}
	if len(result.Suggestions) != 0 {
		t.Fatalf("expected no suggestions, got %v", result.Suggestions)
	}
}

func TestGitRepositoryCheckerStatusError(t *testing.T) {
	requireGit(t)
	tmp := t.TempDir()

	gitDir := filepath.Join(tmp, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	// Remove git from PATH to force status command failure.
	t.Setenv("PATH", tmp)

	c := &GitRepositoryChecker{}
	result := c.Check()

	if result.Status != "ERROR" {
		t.Fatalf("expected ERROR, got %s", result.Status)
	}
}

func TestFilePermissionCheckerCheck(t *testing.T) {
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	restoreConfigState(t)

	c := &FilePermissionChecker{}
	result := c.Check()

	if result.Status != "OK" {
		t.Fatalf("expected OK, got %s", result.Status)
	}
}

func TestFormatDiagnosticResults(t *testing.T) {
	requireGit(t)
	home := t.TempDir()
	setCleanConfigEnvs(t, home)
	writeGitConfig(t, home)
	restoreConfigState(t)

	repo := initRealRepo(t)

	cfgDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatalf("mkdir cfg: %v", err)
	}
	jsonPath := filepath.Join(cfgDir, "config.json")
	if err := os.WriteFile(jsonPath, []byte(`{"writeEmoji": true}`), 0644); err != nil {
		t.Fatalf("write json: %v", err)
	}

	_ = repo

	output := FormatDiagnosticResults()

	if !strings.Contains(output, "System Diagnostics") {
		t.Fatalf("missing system diagnostics section: %s", output)
	}
	if !strings.Contains(output, "Diagnostic Results:") {
		t.Fatalf("missing diagnostic results section: %s", output)
	}
	if !strings.Contains(output, "Diagnostic Summary:") {
		t.Fatalf("missing diagnostic summary section: %s", output)
	}
}
