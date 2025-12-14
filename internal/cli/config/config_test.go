// Tests in this file must not run in parallel because they modify global configuration state.

package config_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"

	cli "github.com/HidemaruOwO/pummit/internal/cli/config"
	cfg "github.com/HidemaruOwO/pummit/internal/config"
)

// setupConfig prepares an isolated config for each test to avoid side effects.
func setupConfig(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, ".config"))
	t.Setenv("LC_ALL", "C")
	t.Setenv("EDITOR", "true")
	prevNoColor := color.NoColor
	color.NoColor = true
	t.Cleanup(func() { color.NoColor = prevNoColor })

	prevCfg := cfg.CurrentTOMLConfig
	cfg.CurrentTOMLConfig = cfg.GetDefaultTOMLConfig()
	t.Cleanup(func() { cfg.CurrentTOMLConfig = prevCfg })

	configPath := filepath.Join(tmp, "config.toml")
	prevPath := cfg.TOMLConfigPath
	cfg.TOMLConfigPath = configPath
	t.Cleanup(func() { cfg.TOMLConfigPath = prevPath })

	if err := cfg.SaveTOMLConfig(); err != nil {
		t.Fatalf("failed to save temp config: %v", err)
	}
	return configPath
}

// loadConfig reads the config file at the given path for assertions.
func loadConfig(t *testing.T, configPath string) cfg.TOMLConfig {
	t.Helper()
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	var c cfg.TOMLConfig
	if err := toml.Unmarshal(data, &c); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	return c
}

// captureOutput captures stdout and stderr produced by f.
func captureOutput(t *testing.T, f func() error) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}

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

func TestConfigListGetSetValidateReset(t *testing.T) {
	configPath := setupConfig(t)

	out, err := captureOutput(t, func() error {
		cli.Cmd.SetArgs([]string{"list"})
		return cli.Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if !strings.Contains(out, "base") {
		t.Fatalf("list output missing base: %s", out)
	}

	out, err = captureOutput(t, func() error {
		cli.Cmd.SetArgs([]string{"get", "base.emoji"})
		return cli.Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if strings.TrimSpace(out) != "true" {
		t.Fatalf("unexpected get output: %s", out)
	}

	cli.Cmd.SetArgs([]string{"set", "base.emoji", "false"})
	if err := cli.Cmd.Execute(); err != nil {
		t.Fatalf("set command failed: %v", err)
	}

	updated := loadConfig(t, configPath)
	if updated.Base.Emoji {
		t.Fatalf("expected base.emoji=false, got true")
	}

	out, err = captureOutput(t, func() error {
		cli.Cmd.SetArgs([]string{"validate"})
		return cli.Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if !strings.Contains(out, "Configuration is valid") {
		t.Fatalf("validate output missing success message: %s", out)
	}

	cfg.CurrentTOMLConfig.Base.Emoji = false
	if err := cfg.SaveTOMLConfig(); err != nil {
		t.Fatalf("failed to save modified config: %v", err)
	}

	_, err = captureOutput(t, func() error {
		cli.Cmd.SetArgs([]string{"reset", "--force"})
		return cli.Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("reset failed: %v", err)
	}

	restored := loadConfig(t, configPath)
	if !reflect.DeepEqual(restored, cfg.GetDefaultTOMLConfig()) {
		t.Fatalf("config not reset to defaults")
	}
}

func TestConfigGetUnknownKey(t *testing.T) {
	setupConfig(t)

	_, err := captureOutput(t, func() error {
		cli.Cmd.SetArgs([]string{"get", "unknown.key"})
		return cli.Cmd.Execute()
	})
	if err == nil {
		t.Fatalf("expected error for unknown key")
	}
}

func TestConfigDefaultShowsHelp(t *testing.T) {
	setupConfig(t)

	out, err := captureOutput(t, func() error {
		cli.Cmd.SetArgs([]string{})
		return cli.Cmd.Execute()
	})
	if err != nil {
		t.Fatalf("default config command failed: %v", err)
	}
	// Should show help containing usage information
	if !strings.Contains(out, "Available subcommands") && !strings.Contains(out, "Usage:") {
		t.Fatalf("expected help output, got: %s", out)
	}
}

func TestValidateFailsOnInvalidConfig(t *testing.T) {
	configPath := setupConfig(t)
	if err := os.WriteFile(configPath, []byte("[meta]\nversion="), 0644); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}

	_, err := captureOutput(t, func() error {
		cli.Cmd.SetArgs([]string{"validate"})
		return cli.Cmd.Execute()
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "failed to parse") {
		t.Fatalf("unexpected error: %v", err)
	}
}
