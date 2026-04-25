package mcp

import (
	"path/filepath"
	"strings"
	"testing"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
)

func TestConfigGetText(t *testing.T) {
	service := usecase.NewConfigService(rootconfig.NewStore(filepath.Join(t.TempDir(), "config.toml")))
	text, err := configGetText(service, "base.emoji")
	if err != nil {
		t.Fatalf("configGetText returned error: %v", err)
	}
	if text != "base.emoji: true" {
		t.Fatalf("configGetText() = %q, want %q", text, "base.emoji: true")
	}

	text, err = configGetText(service, "")
	if err != nil {
		t.Fatalf("configGetText all returned error: %v", err)
	}
	if !strings.Contains(text, "[meta]") {
		t.Fatalf("configGetText all output unexpected: %q", text)
	}
}
