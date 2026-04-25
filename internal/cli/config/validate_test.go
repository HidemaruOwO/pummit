package config_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	rootcli "github.com/HidemaruOwO/pummit/internal/cli"
)

func TestValidateCommand(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("APPDATA", filepath.Join(home, "AppData", "Roaming"))

	configDir := filepath.Join(home, ".config", "pummit")
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(home, "AppData", "Roaming", "pummit")
	}

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte("[base]\nemoji=false\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := rootcli.Execute([]string{"config", "validate"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if stdout.String() != "Configuration is valid\n" {
		t.Fatalf("stdout = %q, want success output", stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}
