// Tests in this file must not run in parallel because they modify global configuration state.

package alias_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"

	aliaspkg "github.com/HidemaruOwO/pummit/internal/alias"
	aliascli "github.com/HidemaruOwO/pummit/internal/cli/alias"
	"github.com/HidemaruOwO/pummit/internal/config"
)

// setupConfig prepares an isolated config for each test to avoid side effects.
func setupConfig(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, ".config"))
	t.Setenv("LC_ALL", "C")
	prevNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() { color.NoColor = prevNoColor })

	prevCfg := config.CurrentTOMLConfig
	config.CurrentTOMLConfig = config.GetDefaultTOMLConfig()
	t.Cleanup(func() { config.CurrentTOMLConfig = prevCfg })

	configPath := filepath.Join(tmp, "config.toml")
	prevPath := config.TOMLConfigPath
	config.TOMLConfigPath = configPath
	t.Cleanup(func() { config.TOMLConfigPath = prevPath })

	if err := config.SaveTOMLConfig(); err != nil {
		t.Fatalf("failed to save temp config: %v", err)
	}
	return configPath
}

// loadConfig reads the config file at the given path for assertions.
func loadConfig(t *testing.T, configPath string) config.TOMLConfig {
	t.Helper()
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	var cfg config.TOMLConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	return cfg
}

// captureOutput captures stdout and stderr produced by f.
func captureOutput(t *testing.T, f func() error) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	defer func() {
		if err := w.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			t.Fatalf("failed to close writer: %v", err)
		}
	}()

	oldStdout, oldStderr := os.Stdout, os.Stderr
	oldColorOut, oldColorErr := color.Output, color.Error
	os.Stdout, os.Stderr = w, w
	color.Output, color.Error = w, w
	defer func() {
		os.Stdout, os.Stderr = oldStdout, oldStderr
		color.Output, color.Error = oldColorOut, oldColorErr
	}()

	runErr := f()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("failed to close reader: %v", err)
	}

	return buf.String(), runErr
}

// hasShortcut checks whether shortcut exists in cfg.
func hasShortcut(cfg config.TOMLConfig, sc string) bool {
	for _, e := range cfg.Alias.Entries {
		for _, s := range e.Shortcuts {
			if s == sc {
				return true
			}
		}
	}
	return false
}

func TestAliasLifecycle(t *testing.T) {
	configPath := setupConfig(t)

	aliascli.AddCmd.SetArgs([]string{"fs", "sparkles"})
	if err := aliascli.AddCmd.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	cfg := loadConfig(t, configPath)
	if !hasShortcut(cfg, "fs") {
		t.Fatalf("alias not added: %+v", cfg.Alias.Entries)
	}

	out, err := captureOutput(t, func() error {
		aliascli.ListCmd.SetArgs(nil)
		return aliascli.ListCmd.Execute()
	})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "fs") {
		t.Fatalf("list output missing alias: %s", out)
	}

	aliascli.DeleteCmd.SetArgs([]string{"fs", "--confirm"})
	if err := aliascli.DeleteCmd.Execute(); err != nil {
		t.Fatalf("delete command failed: %v", err)
	}

	cfg = loadConfig(t, configPath)
	if hasShortcut(cfg, "fs") {
		t.Fatalf("alias not deleted: %+v", cfg.Alias.Entries)
	}

	aliascli.AddCmd.SetArgs([]string{"fs", "sparkles"})
	if err := aliascli.AddCmd.Execute(); err != nil {
		t.Fatalf("add command (second time) failed: %v", err)
	}

	aliascli.ResetCmd.SetArgs([]string{"--confirm"})
	if err := aliascli.ResetCmd.Execute(); err != nil {
		t.Fatalf("reset command failed: %v", err)
	}

	cfg = loadConfig(t, configPath)
	if hasShortcut(cfg, "fs") {
		t.Fatalf("alias not reset: %+v", cfg.Alias.Entries)
	}
	if !reflect.DeepEqual(cfg.Alias, config.GetDefaultTOMLConfig().Alias) {
		t.Fatalf("aliases not restored to defaults")
	}
}

func TestAliasAddDuplicate(t *testing.T) {
	configPath := setupConfig(t)

	aliascli.AddCmd.SetArgs([]string{"fs", "sparkles"})
	if err := aliascli.AddCmd.Execute(); err != nil {
		t.Fatalf("initial add failed: %v", err)
	}

	out, err := captureOutput(t, func() error {
		aliascli.AddCmd.SetArgs([]string{"fs", "sparkles"})
		return aliascli.AddCmd.Execute()
	})
	if !errors.Is(err, aliaspkg.ErrAliasExists) {
		t.Fatalf("expected ErrAliasExists, got %v", err)
	}
	if out == "" {
		t.Fatalf("expected error output")
	}

	cfg := loadConfig(t, configPath)
	count := 0
	for _, e := range cfg.Alias.Entries {
		for _, sc := range e.Shortcuts {
			if sc == "fs" {
				count++
			}
		}
	}
	if count != 1 {
		t.Fatalf("duplicate alias present: %+v", cfg.Alias.Entries)
	}
}

func TestAliasDeleteNonExistent(t *testing.T) {
	configPath := setupConfig(t)

	out, err := captureOutput(t, func() error {
		aliascli.DeleteCmd.SetArgs([]string{"unknown", "--confirm"})
		return aliascli.DeleteCmd.Execute()
	})
	if !errors.Is(err, aliaspkg.ErrAliasNotFound) {
		t.Fatalf("expected ErrAliasNotFound, got %v", err)
	}
	if out == "" {
		t.Fatalf("expected error output")
	}

	cfg := loadConfig(t, configPath)
	if !reflect.DeepEqual(cfg.Alias, config.GetDefaultTOMLConfig().Alias) {
		t.Fatalf("config changed after failed delete")
	}
}
