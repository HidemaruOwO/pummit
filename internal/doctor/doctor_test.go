package doctor

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/HidemaruOwO/pummit/internal/config"
)

// prependToPATH inserts dir at the beginning of PATH to control git resolution.
func prependToPATH(t *testing.T, dir string) {
	t.Helper()
	sep := string(os.PathListSeparator)
	old := os.Getenv("PATH")
	t.Setenv("PATH", dir+sep+old)
}

// setCleanConfigEnvs isolates global config state for deterministic tests.
func setCleanConfigEnvs(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
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

// skipOnWindows avoids shell script execution on unsupported platforms.
func skipOnWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell script fake git not supported on Windows")
	}
}

// addFakeGit writes a minimal git mock used to simulate command responses.
func addFakeGit(t *testing.T, dir string) {
	t.Helper()
	script := `#!/bin/sh
set -eu
case "$1" in
  --version)
    echo 'git version 2.39.0'
    ;;
  config)
    if [ "$2" = '--global' ] && [ "$3" = 'user.name' ]; then
      echo 'Test User'
    elif [ "$2" = '--global' ] && [ "$3" = 'user.email' ]; then
      echo 'test@example.com'
    else
      exit 1
    fi
    ;;
  status)
    echo ''
    ;;
  *)
    exit 1
    ;;
esac
`
	path := filepath.Join(dir, "git")
	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		t.Fatalf("write fake git: %v", err)
	}
	prependToPATH(t, dir)
}

func TestRunAllChecksHealthyEnvNoErrors(t *testing.T) {
	skipOnWindows(t)

	home := t.TempDir()
	setCleanConfigEnvs(t, home)

	configDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	prevCfg := config.CurrentTOMLConfig
	prevPath := config.TOMLConfigPath
	t.Cleanup(func() {
		config.CurrentTOMLConfig = prevCfg
		config.TOMLConfigPath = prevPath
	})
	config.CurrentTOMLConfig = config.GetDefaultTOMLConfig()
	config.TOMLConfigPath = filepath.Join(configDir, "config.toml")
	if err := config.SaveTOMLConfig(); err != nil {
		t.Fatalf("save config: %v", err)
	}

	fakeBin := filepath.Join(home, "bin")
	if err := os.MkdirAll(fakeBin, 0755); err != nil {
		t.Fatalf("mkdir fake bin: %v", err)
	}
	addFakeGit(t, fakeBin)

	repoDir := filepath.Join(home, "repo")
	if err := os.MkdirAll(filepath.Join(repoDir, ".git"), 0755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	results := RunAllChecks()
	if results.HasErrors() {
		t.Fatalf("expected no errors, got %+v", results)
	}
}

func TestRunAllChecksGitMissing(t *testing.T) {
	skipOnWindows(t)

	home := t.TempDir()
	setCleanConfigEnvs(t, home)

	configDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	prevCfg := config.CurrentTOMLConfig
	prevPath := config.TOMLConfigPath
	t.Cleanup(func() {
		config.CurrentTOMLConfig = prevCfg
		config.TOMLConfigPath = prevPath
	})
	config.CurrentTOMLConfig = config.GetDefaultTOMLConfig()
	config.TOMLConfigPath = filepath.Join(configDir, "config.toml")
	if err := config.SaveTOMLConfig(); err != nil {
		t.Fatalf("save config: %v", err)
	}

	repoDir := filepath.Join(home, "repo")
	if err := os.MkdirAll(filepath.Join(repoDir, ".git"), 0755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	t.Setenv("PATH", home)

	results := RunAllChecks()
	var gitErr bool
	for _, r := range results {
		if r.Name == "Git Configuration" && r.Status == "ERROR" {
			gitErr = true
		}
	}
	if !gitErr {
		t.Fatalf("expected git missing error, got %+v", results)
	}
}

func TestRunAllChecksMalformedConfig(t *testing.T) {
	skipOnWindows(t)

	home := t.TempDir()
	setCleanConfigEnvs(t, home)

	configDir := filepath.Join(home, ".config", "pummit")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("mkdir config: %v", err)
	}
	bad := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(bad, []byte("[[invalid"), 0644); err != nil {
		t.Fatalf("write bad config: %v", err)
	}
	prevCfg := config.CurrentTOMLConfig
	prevPath := config.TOMLConfigPath
	t.Cleanup(func() {
		config.CurrentTOMLConfig = prevCfg
		config.TOMLConfigPath = prevPath
	})

	fakeBin := filepath.Join(home, "bin")
	if err := os.MkdirAll(fakeBin, 0755); err != nil {
		t.Fatalf("mkdir fake bin: %v", err)
	}
	addFakeGit(t, fakeBin)

	repoDir := filepath.Join(home, "repo")
	if err := os.MkdirAll(filepath.Join(repoDir, ".git"), 0755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repoDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	results := RunAllChecks()
	var cfgErr bool
	for _, r := range results {
		if r.Name == "Configuration Files" && r.Status == "ERROR" &&
			strings.Contains(r.Message, "TOML configuration file is corrupted") {
			cfgErr = true
		}
	}
	if !cfgErr {
		t.Fatalf("expected config parse error, got %+v", results)
	}
}
