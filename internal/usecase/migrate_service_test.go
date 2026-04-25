package usecase

import (
	"os"
	"path/filepath"
	"testing"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/variable"
)

func TestMigrateServiceLifecycle(t *testing.T) {
	dir := t.TempDir()
	store := rootconfig.NewStore(filepath.Join(dir, "config.toml"))
	service := NewMigrateService(store)
	jsonPath, err := store.JSONPath()
	if err != nil {
		t.Fatalf("JSONPath returned error: %v", err)
	}
	if err := os.WriteFile(jsonPath, []byte(variable.DEFAULT_CONFIG), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	status, err := service.Status()
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if status.Status != MigrationStatusJSONOnly {
		t.Fatalf("status = %s, want %s", status.Status, MigrationStatusJSONOnly)
	}

	result, err := service.Migrate(false)
	if err != nil {
		t.Fatalf("Migrate returned error: %v", err)
	}
	if result.Message == "" {
		t.Fatal("Migrate returned empty message")
	}
	if _, err := os.Stat(filepath.Join(dir, "config.toml")); err != nil {
		t.Fatalf("config.toml missing after migrate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.json.bak")); err != nil {
		t.Fatalf("backup missing after migrate: %v", err)
	}

	result, err = service.Rollback(true)
	if err != nil {
		t.Fatalf("Rollback returned error: %v", err)
	}
	if result.Message == "" {
		t.Fatal("Rollback returned empty message")
	}
	if _, err := os.Stat(jsonPath); err != nil {
		t.Fatalf("config.json missing after rollback: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.toml")); !os.IsNotExist(err) {
		t.Fatalf("config.toml should be removed after rollback")
	}
}
