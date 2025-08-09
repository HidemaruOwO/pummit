// NOTE: These tests mutate package-level globals.
// Do NOT use t.Parallel() in this file, as doing so may cause tests to interfere with
// each other and produce unreliable results. t.Parallel() is only safe if tests do not
// mutate shared state (such as package-level variables).
package config

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

// isolateEnv prepares HOME for an isolated test and resets globals.
func isolateEnv(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, ".config"))

	prevCfg := CurrentTOMLConfig
	prevPath := TOMLConfigPath
	CurrentTOMLConfig = TOMLConfig{}
	TOMLConfigPath = ""

	t.Cleanup(func() {
		CurrentTOMLConfig = prevCfg
		TOMLConfigPath = prevPath
	})
}

// readConfigFile loads TOMLConfig from disk for verification.
func readConfigFile(t *testing.T) TOMLConfig {
	t.Helper()
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	p := filepath.Join(dir, "config.toml")
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("Stat: %v", err)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var cfg TOMLConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	return cfg
}

// TestTOMLConfigRoundTrip ensures saving and loading are idempotent.
func TestTOMLConfigRoundTrip(t *testing.T) {
	isolateEnv(t)
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	if err := SaveTOMLConfig(); err != nil {
		t.Fatalf("SaveTOMLConfig: %v", err)
	}
	CurrentTOMLConfig = TOMLConfig{}
	TOMLConfigPath = ""
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	want := GetDefaultTOMLConfig()
	if !reflect.DeepEqual(want, CurrentTOMLConfig) {
		t.Fatalf("CurrentTOMLConfig = %#v want %#v", CurrentTOMLConfig, want)
	}
	fileCfg := readConfigFile(t)
	if !reflect.DeepEqual(want, fileCfg) {
		t.Fatalf("fileCfg = %#v want %#v", fileCfg, want)
	}
}

// TestLoadTOMLConfigCreatesDefault verifies a default file is generated.
func TestLoadTOMLConfigCreatesDefault(t *testing.T) {
	isolateEnv(t)
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	want := GetDefaultTOMLConfig()
	if !reflect.DeepEqual(want, CurrentTOMLConfig) {
		t.Fatalf("CurrentTOMLConfig = %#v want %#v", CurrentTOMLConfig, want)
	}
	_ = readConfigFile(t)
}

// TestLoadTOMLConfigInvalid checks invalid TOML handling.
func TestLoadTOMLConfigInvalid(t *testing.T) {
	isolateEnv(t)
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte("::invalid::"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := LoadTOMLConfig(); err == nil {
		t.Fatalf("expected error for invalid TOML")
	} else {
		var parseErr toml.ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("unexpected error type: %v", err)
		}
	}
}

// TestSaveTOMLConfigError covers directory creation failures.
func TestSaveTOMLConfigError(t *testing.T) {
	isolateEnv(t)
	home := os.Getenv("HOME")
	if err := os.WriteFile(filepath.Join(home, ".config"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	if err := SaveTOMLConfig(); err == nil {
		t.Fatalf("expected error when config dir is file")
	}
}

// TestLoadTOMLConfigCreateError covers failure when default file cannot be created.
func TestLoadTOMLConfigCreateError(t *testing.T) {
	isolateEnv(t)
	home := os.Getenv("HOME")
	if err := os.WriteFile(filepath.Join(home, ".config"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := LoadTOMLConfig(); err == nil {
		t.Fatalf("expected error when config dir is file")
	}
}

// TestLoadTOMLConfigUnknownKeys ensures unknown keys are ignored.
func TestLoadTOMLConfigUnknownKeys(t *testing.T) {
	isolateEnv(t)
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(GetDefaultTOMLConfig()); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	buf.WriteString("unknown = 1\n")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	want := GetDefaultTOMLConfig()
	if !reflect.DeepEqual(want, CurrentTOMLConfig) {
		t.Fatalf("CurrentTOMLConfig = %#v want %#v", CurrentTOMLConfig, want)
	}
}

// TestLoadTOMLConfigPartial verifies missing fields are zero values.
func TestLoadTOMLConfigPartial(t *testing.T) {
	isolateEnv(t)
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	content := []byte("[base]\nemoji=false\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	// The design intentionally leaves unspecified fields at zero values
	// instead of merging with defaults.
	if CurrentTOMLConfig.Base.Emoji != false {
		t.Fatalf("Base.Emoji got %v", CurrentTOMLConfig.Base.Emoji)
	}
	if CurrentTOMLConfig.Base.FilesLength != 0 {
		t.Fatalf("Base.FilesLength want 0 got %d", CurrentTOMLConfig.Base.FilesLength)
	}
}

// TestSaveTOMLConfigPermissionError checks write permission failures.
func TestSaveTOMLConfigPermissionError(t *testing.T) {
	if !isUnixNonRoot() {
		t.Skip("skipping permission test: requires non-root Unix-like system")
	}
	isolateEnv(t)
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatalf("Chmod(file): %v", err)
	}
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatalf("Chmod(dir): %v", err)
	}
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	if err := SaveTOMLConfig(); err == nil {
		t.Fatalf("expected permission error")
	}
}

// TestTOMLConfigPathOverride documents path override behavior.
func TestTOMLConfigPathOverride(t *testing.T) {
	isolateEnv(t)
	custom := filepath.Join(t.TempDir(), "custom.toml")
	TOMLConfigPath = custom
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	if err := SaveTOMLConfig(); err != nil {
		t.Fatalf("SaveTOMLConfig: %v", err)
	}
	if _, err := os.Stat(custom); err != nil {
		t.Fatalf("Stat: %v", err)
	}
	CurrentTOMLConfig = TOMLConfig{}
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	// Specification A: explicit override path remains after Load.
	if TOMLConfigPath != custom {
		t.Fatalf("TOMLConfigPath = %s want %s", TOMLConfigPath, custom)
	}
}
