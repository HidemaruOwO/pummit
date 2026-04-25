package mcp

import (
	"testing"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
)

func TestFormatAliasEntries(t *testing.T) {
	entries := []domainconfig.AliasEntry{{Shortcuts: []string{"s", "feat"}, Name: "sparkles", Emoji: "✨"}}
	got := formatAliasEntries(entries)
	want := "✨ sparkles => s,feat"
	if got != want {
		t.Fatalf("formatAliasEntries() = %q, want %q", got, want)
	}
}
