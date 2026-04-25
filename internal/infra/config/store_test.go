package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
)

func TestConfigDir(t *testing.T) {
	store := NewStore("")
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("APPDATA", filepath.Join(home, "AppData", "Roaming"))

	dir, err := store.ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir returned error: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		want := filepath.Join(home, "AppData", "Roaming", "pummit")
		if dir != want {
			t.Fatalf("ConfigDir = %q, want %q", dir, want)
		}
	default:
		want := filepath.Join(home, ".config", "pummit")
		if dir != want {
			t.Fatalf("ConfigDir = %q, want %q", dir, want)
		}
	}
}

func TestLoadCreatesDefaultConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	store := NewStore(path)

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	defaults := domainconfig.Default()
	if cfg.Meta.Version != defaults.Meta.Version {
		t.Fatalf("Meta.Version = %q, want %q", cfg.Meta.Version, defaults.Meta.Version)
	}

	if cfg.Base.FilesLength != defaults.Base.FilesLength {
		t.Fatalf("Base.FilesLength = %d, want %d", cfg.Base.FilesLength, defaults.Base.FilesLength)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config file not created: %v", err)
	}
}

func TestLoadMergesPartialConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("[base]\nemoji=false\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	store := NewStore(path)
	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Base.Emoji {
		t.Fatalf("Base.Emoji = true, want false")
	}

	if cfg.Base.FilesLength != 50 {
		t.Fatalf("Base.FilesLength = %d, want 50", cfg.Base.FilesLength)
	}

	if cfg.Locale.Language != "ja" {
		t.Fatalf("Locale.Language = %q, want ja", cfg.Locale.Language)
	}
}

func TestLoadRejectsUnknownKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte("unknown = 1\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	store := NewStore(path)
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for unknown keys")
	}

	if !strings.Contains(err.Error(), "unknown keys") {
		t.Fatalf("error = %q, want unknown keys", err.Error())
	}
}

func TestLoadMigratesLegacyJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	jsonPath := filepath.Join(dir, "config.json")
	data := []byte(`{"writeEmoji":false,"useAlias":true,"useLimitPathesLength":false,"alias":[["x,y","sparkles","✨"]]}`)
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	store := NewStore(path)
	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Base.Emoji {
		t.Fatalf("Base.Emoji = true, want false")
	}

	if cfg.Base.FilesLength != 0 {
		t.Fatalf("Base.FilesLength = %d, want 0", cfg.Base.FilesLength)
	}

	if len(cfg.Alias.Entries) != 1 || len(cfg.Alias.Entries[0].Shortcuts) != 2 {
		t.Fatalf("unexpected migrated aliases: %+v", cfg.Alias.Entries)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config.toml not written: %v", err)
	}
}
