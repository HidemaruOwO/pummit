package alias

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/HidemaruOwO/pummit/internal/config"
)

// setupTestConfig creates an isolated TOML config path and resets the in-memory config.
func setupTestConfig(t *testing.T) func() {
	t.Helper()

	origPath := config.TOMLConfigPath
	origConfig := config.CurrentTOMLConfig

	tmpDir := t.TempDir()
	config.TOMLConfigPath = filepath.Join(tmpDir, "config.toml")
	config.CurrentTOMLConfig = config.TOMLConfig{
		Alias: config.AliasConfig{
			Enabled: true,
			Entries: []config.AliasEntry{},
		},
	}

	return func() {
		config.TOMLConfigPath = origPath
		config.CurrentTOMLConfig = origConfig
	}
}

func cloneEntries(entries []config.AliasEntry) []config.AliasEntry {
	cloned := make([]config.AliasEntry, len(entries))
	for i, e := range entries {
		shortcuts := make([]string, len(e.Shortcuts))
		copy(shortcuts, e.Shortcuts)
		cloned[i] = config.AliasEntry{
			Shortcuts: shortcuts,
			Name:      e.Name,
			Emoji:     e.Emoji,
		}
	}
	return cloned
}

func TestGetEmoji(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	tests := []struct {
		name      string
		entries   []config.AliasEntry
		lookup    string
		wantFound bool
		wantName  string
		wantEmoji string
	}{
		{
			name: "found",
			entries: []config.AliasEntry{
				{Shortcuts: []string{"s", "feat"}, Name: "sparkles", Emoji: "✨"},
			},
			lookup:    "feat",
			wantFound: true,
			wantName:  "sparkles",
			wantEmoji: "✨",
		},
		{
			name: "not found",
			entries: []config.AliasEntry{
				{Shortcuts: []string{"s"}, Name: "sparkles", Emoji: "✨"},
			},
			lookup:    "unknown",
			wantFound: false,
			wantName:  "",
			wantEmoji: "",
		},
		{
			name:      "empty entries",
			entries:   []config.AliasEntry{},
			lookup:    "any",
			wantFound: false,
			wantName:  "",
			wantEmoji: "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentTOMLConfig.Alias.Entries = cloneEntries(tt.entries)

			found, gotName, gotEmoji := GetEmoji(tt.lookup)

			if found != tt.wantFound {
				t.Fatalf("found mismatch: got %v, want %v", found, tt.wantFound)
			}
			if gotName != tt.wantName || gotEmoji != tt.wantEmoji {
				t.Fatalf("unexpected result: name=%q emoji=%q, want name=%q emoji=%q", gotName, gotEmoji, tt.wantName, tt.wantEmoji)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	tests := []struct {
		name           string
		initialEntries []config.AliasEntry
		shortcut       string
		aliasName      string
		emoji          string
		wantErr        error
		assert         func(t *testing.T)
	}{
		{
			name:           "add new alias",
			initialEntries: []config.AliasEntry{},
			shortcut:       "n",
			aliasName:      "new",
			emoji:          "🆕",
			wantErr:        nil,
			assert: func(t *testing.T) {
				_, entry := findAlias("n")
				if entry == nil {
					t.Fatalf("alias not added")
				}
				if entry.Name != "new" || entry.Emoji != "🆕" {
					t.Fatalf("unexpected entry: name=%s emoji=%s", entry.Name, entry.Emoji)
				}
				if len(entry.Shortcuts) != 1 || entry.Shortcuts[0] != "n" {
					t.Fatalf("unexpected shortcuts: %v", entry.Shortcuts)
				}
			},
		},
		{
			name: "shortcut already exists",
			initialEntries: []config.AliasEntry{
				{Shortcuts: []string{"s"}, Name: "sparkles", Emoji: "✨"},
			},
			shortcut:  "s",
			aliasName: "sparkles",
			emoji:     "✨",
			wantErr:   ErrAliasExists,
			assert: func(t *testing.T) {
				if len(config.CurrentTOMLConfig.Alias.Entries) != 1 {
					t.Fatalf("entries should remain unchanged")
				}
			},
		},
		{
			name: "append shortcut to existing emoji with same name",
			initialEntries: []config.AliasEntry{
				{Shortcuts: []string{"s"}, Name: "sparkles", Emoji: "✨"},
			},
			shortcut:  "feat",
			aliasName: "sparkles",
			emoji:     "✨",
			wantErr:   nil,
			assert: func(t *testing.T) {
				_, entry := findAlias("feat")
				if entry == nil {
					t.Fatalf("alias not added to existing emoji")
				}
				if len(entry.Shortcuts) != 2 {
					t.Fatalf("expected 2 shortcuts, got %d", len(entry.Shortcuts))
				}
			},
		},
		{
			name: "emoji exists with different name",
			initialEntries: []config.AliasEntry{
				{Shortcuts: []string{"s"}, Name: "sparkles", Emoji: "✨"},
			},
			shortcut:  "c",
			aliasName: "construction",
			emoji:     "✨",
			wantErr:   errors.New("different name"),
			assert: func(t *testing.T) {
				if len(config.CurrentTOMLConfig.Alias.Entries) != 1 {
					t.Fatalf("entries should remain unchanged on error")
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentTOMLConfig.Alias.Entries = cloneEntries(tt.initialEntries)

			err := Add(tt.shortcut, tt.aliasName, tt.emoji)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Add returned unexpected error: %v", err)
				}
			} else {
				if !errors.Is(err, tt.wantErr) && err == nil {
					t.Fatalf("expected error but got nil")
				}
				if errors.Is(tt.wantErr, ErrAliasExists) && !errors.Is(err, ErrAliasExists) {
					t.Fatalf("expected ErrAliasExists, got %v", err)
				}
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
			}

			if tt.assert != nil {
				tt.assert(t)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	tests := []struct {
		name           string
		initialEntries []config.AliasEntry
		shortcut       string
		wantErr        error
		assert         func(t *testing.T)
	}{
		{
			name: "delete existing single shortcut",
			initialEntries: []config.AliasEntry{
				{Shortcuts: []string{"s"}, Name: "sparkles", Emoji: "✨"},
			},
			shortcut: "s",
			wantErr:  nil,
			assert: func(t *testing.T) {
				if len(config.CurrentTOMLConfig.Alias.Entries) != 0 {
					t.Fatalf("expected entries to be empty after deletion")
				}
			},
		},
		{
			name: "shortcut not found",
			initialEntries: []config.AliasEntry{
				{Shortcuts: []string{"s"}, Name: "sparkles", Emoji: "✨"},
			},
			shortcut: "x",
			wantErr:  ErrAliasNotFound,
			assert: func(t *testing.T) {
				if len(config.CurrentTOMLConfig.Alias.Entries) != 1 {
					t.Fatalf("entries should remain unchanged")
				}
			},
		},
		{
			name: "delete last shortcut removes entry only",
			initialEntries: []config.AliasEntry{
				{Shortcuts: []string{"a"}, Name: "art", Emoji: "🎨"},
				{Shortcuts: []string{"b"}, Name: "bug", Emoji: "🐛"},
			},
			shortcut: "b",
			wantErr:  nil,
			assert: func(t *testing.T) {
				if len(config.CurrentTOMLConfig.Alias.Entries) != 1 {
					t.Fatalf("expected one entry remaining")
				}
				if _, entry := findAlias("a"); entry == nil {
					t.Fatalf("remaining entry missing")
				}
			},
		},
		{
			name: "delete one of multiple shortcuts",
			initialEntries: []config.AliasEntry{
				{Shortcuts: []string{"x", "y"}, Name: "sparkles", Emoji: "✨"},
			},
			shortcut: "y",
			wantErr:  nil,
			assert: func(t *testing.T) {
				_, entry := findAlias("x")
				if entry == nil {
					t.Fatalf("entry should remain after deleting one shortcut")
				}
				if len(entry.Shortcuts) != 1 || entry.Shortcuts[0] != "x" {
					t.Fatalf("unexpected shortcuts after delete: %v", entry.Shortcuts)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentTOMLConfig.Alias.Entries = cloneEntries(tt.initialEntries)

			err := Delete(tt.shortcut)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Delete returned unexpected error: %v", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.assert != nil {
				tt.assert(t)
			}
		})
	}
}

func TestGet(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	tests := []struct {
		name      string
		entries   []config.AliasEntry
		lookup    string
		wantEmoji string
		wantErr   error
	}{
		{
			name: "existing shortcut",
			entries: []config.AliasEntry{
				{Shortcuts: []string{"s"}, Name: "sparkles", Emoji: "✨"},
			},
			lookup:    "s",
			wantEmoji: "✨",
			wantErr:   nil,
		},
		{
			name:      "missing shortcut",
			entries:   []config.AliasEntry{},
			lookup:    "unknown",
			wantEmoji: "",
			wantErr:   ErrAliasNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentTOMLConfig.Alias.Entries = cloneEntries(tt.entries)

			gotEmoji, err := Get(tt.lookup)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Get returned unexpected error: %v", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if gotEmoji != tt.wantEmoji {
				t.Fatalf("emoji mismatch: got %q, want %q", gotEmoji, tt.wantEmoji)
			}
		})
	}
}

func TestList(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	tests := []struct {
		name    string
		entries []config.AliasEntry
		want    map[string]string
	}{
		{
			name: "multiple entries",
			entries: []config.AliasEntry{
				{Shortcuts: []string{"s", "feat"}, Name: "sparkles", Emoji: "✨"},
				{Shortcuts: []string{"b", "bug"}, Name: "bug", Emoji: "🐛"},
			},
			want: map[string]string{
				"s":    "✨",
				"feat": "✨",
				"b":    "🐛",
				"bug":  "🐛",
			},
		},
		{
			name:    "empty",
			entries: []config.AliasEntry{},
			want:    map[string]string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config.CurrentTOMLConfig.Alias.Entries = cloneEntries(tt.entries)

			got := List()

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("list mismatch: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReset(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	config.CurrentTOMLConfig.Alias.Entries = []config.AliasEntry{
		{Shortcuts: []string{"x"}, Name: "custom", Emoji: "❌"},
	}

	if err := Reset(); err != nil {
		t.Fatalf("Reset returned error: %v", err)
	}

	defaultAlias := config.GetDefaultTOMLConfig().Alias
	if !reflect.DeepEqual(config.CurrentTOMLConfig.Alias, defaultAlias) {
		t.Fatalf("alias config not reset to default")
	}
}
