package doctor

import (
  "path/filepath"
  "regexp"
  "runtime"
  "testing"
)

var reGitVersion = regexp.MustCompile(`^\d+\.\d+(?:\.\d+)?$`)

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

func TestGetSystemInfoGitPresentVersionFormat(t *testing.T) {
  tmp := t.TempDir()
  setCleanConfigEnvs(t, tmp)

  info := GetSystemInfo()
  const gitNotFound = "Not installed or not accessible"
  if info.GitVersion == gitNotFound {
    t.Skip("git not installed; skipping")
  }
  if !reGitVersion.MatchString(info.GitVersion) {
    t.Fatalf("unexpected git version format: %q", info.GitVersion)
  }
}

func TestGetSystemInfoGitMissingReported(t *testing.T) {
  tmp := t.TempDir()
  setCleanConfigEnvs(t, tmp)
  t.Setenv("PATH", tmp)

  info := GetSystemInfo()
  const want = "Not installed or not accessible"
  if info.GitVersion != want {
    t.Fatalf("git missing: want %q, got %q", want, info.GitVersion)
  }
}

