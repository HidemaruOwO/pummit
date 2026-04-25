package usecase

import (
	"path/filepath"
	"testing"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	rootemoji "github.com/HidemaruOwO/pummit/internal/infra/emoji"
)

func TestAliasServiceLifecycle(t *testing.T) {
	store := rootconfig.NewStore(filepath.Join(t.TempDir(), "config.toml"))
	configService := NewConfigService(store)
	service := NewAliasService(configService, rootemoji.NewCatalog(), nil)

	if err := service.Add("rs", "rocket", "🚀"); err != nil {
		t.Fatalf("Add returned error: %v", err)
	}

	entries, err := service.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if !hasShortcut(entries, "rs") {
		t.Fatalf("added shortcut not found in entries: %+v", entries)
	}

	if err := service.Delete("rs"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	entries, err = service.List()
	if err != nil {
		t.Fatalf("List after delete returned error: %v", err)
	}
	if hasShortcut(entries, "rs") {
		t.Fatalf("deleted shortcut still present: %+v", entries)
	}

	if err := service.Add("zz", "rocket", "🚀"); err != nil {
		t.Fatalf("second Add returned error: %v", err)
	}
	if err := service.Reset(); err != nil {
		t.Fatalf("Reset returned error: %v", err)
	}
	entries, err = service.List()
	if err != nil {
		t.Fatalf("List after reset returned error: %v", err)
	}
	if hasShortcut(entries, "zz") {
		t.Fatalf("reset shortcut still present: %+v", entries)
	}
}

func hasShortcut(entries []domainconfig.AliasEntry, shortcut string) bool {
	for _, entry := range entries {
		for _, existing := range entry.Shortcuts {
			if existing == shortcut {
				return true
			}
		}
	}
	return false
}
