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
func setupConfig(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, ".config"))
	t.Setenv("LC_ALL", "C")
	color.NoColor = true

	config.CurrentTOMLConfig = config.GetDefaultTOMLConfig()
	config.TOMLConfigPath = ""
	if err := config.SaveTOMLConfig(); err != nil {
		t.Fatalf("failed to save temp config: %v", err)
	}
}

// loadConfig reads the current config file for assertions.
func loadConfig(t *testing.T) config.TOMLConfig {
	t.Helper()
	data, err := os.ReadFile(config.TOMLConfigPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	var cfg config.TOMLConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	return cfg
}

// captureOutput captures stdout produced by f.
func captureOutput(t *testing.T, f func() error) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}

	oldStdout := os.Stdout
	oldColor := color.Output
	os.Stdout = w
	color.Output = w
	defer func() {
		os.Stdout = oldStdout
		color.Output = oldColor
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
	setupConfig(t)

	aliascli.AddCmd.SetArgs([]string{"fs", "sparkles"})
	if err := aliascli.AddCmd.Execute(); err != nil {
		t.Fatalf("add command failed: %v", err)
	}

	cfg := loadConfig(t)
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

	cfg = loadConfig(t)
	if hasShortcut(cfg, "fs") {
		t.Fatalf("alias not deleted: %+v", cfg.Alias.Entries)
	}

	aliascli.AddCmd.SetArgs([]string{"fs", "sparkles"})
	_ = aliascli.AddCmd.Execute()

	aliascli.ResetCmd.SetArgs([]string{"--confirm"})
	if err := aliascli.ResetCmd.Execute(); err != nil {
		t.Fatalf("reset command failed: %v", err)
	}

	cfg = loadConfig(t)
	if hasShortcut(cfg, "fs") {
		t.Fatalf("alias not reset: %+v", cfg.Alias.Entries)
	}
	if !reflect.DeepEqual(cfg.Alias, config.GetDefaultTOMLConfig().Alias) {
		t.Fatalf("aliases not restored to defaults")
	}
}

func TestAliasAddDuplicate(t *testing.T) {
	setupConfig(t)

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

	cfg := loadConfig(t)
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
	setupConfig(t)

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

	cfg := loadConfig(t)
	if !reflect.DeepEqual(cfg.Alias, config.GetDefaultTOMLConfig().Alias) {
		t.Fatalf("config changed after failed delete")
	}
}
