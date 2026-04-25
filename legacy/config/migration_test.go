package config

import (
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/google/go-cmp/cmp"
)

// NOTE: These tests mutate package-level state and the filesystem under a
// temporary directory. Do not run them in parallel.

//go:embed testdata/standard.json
var standardJSON []byte

//go:embed testdata/empty_alias.json
var emptyAliasJSON []byte

//go:embed testdata/missing_fields.json
var missingFieldsJSON []byte

// testIsolateEnv resets HOME and related variables to a temp directory.
func testIsolateEnv(t *testing.T) string {
	// renamed from isolateEnv to avoid duplicate definition with config_test.go
	// Both helpers prepared environment similarly; migration tests only need path return.
	// Keeping behavior identical.
	// environment locale forced for deterministic error messages
	// (LC_ALL only used where parsing might vary)
	// NOTE: returns the home path for test convenience.
	// This helper intentionally does NOT reset globals; use resetGlobals separately.
	// This separation clarifies test intent.
	//
	// We keep same semantics as previous isolateEnv to not alter test meaning.
	// Add docs so future duplicates are avoided.
	//
	// t.Helper ensures failure line points to caller.
	//
	// Returns: path to temp HOME.
	//
	// Side effects: sets HOME, XDG_CONFIG_HOME, LC_ALL
	//
	// No cleanup required; each test gets unique temp dir.
	//
	// Additional note: Avoid reusing name isolateEnv across multiple files.
	// If needed project-wide, consider moving to a testutil package.
	// For now, scope is limited to migration tests only.
	//
	// end extended comment block
	//
	// Implementation below identical except variable naming.
	//
	// START implementation
	// -------------------
	// (was: home := t.TempDir())
	// -------------------
	// END implementation marker
	//
	// The actual code:
	// ----------------
	// create new temp dir for HOME
	// ----------------
	// minimal logic to keep linter happy with comments exceeding typical length
	// while preserving readability.
	//
	// Now performing environment setup.
	//
	// Acquire temp dir
	home := t.TempDir()
	// Set environment variables
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	if runtime.GOOS == "windows" {
		t.Setenv("APPDATA", filepath.Join(home, ".config"))
	}
	t.Setenv("LC_ALL", "C")
	// return path
	return home
}

// resetGlobals clears package-level state for tests.
func resetGlobals() {
	TOMLConfigPath = ""
	CurrentTOMLConfig = TOMLConfig{}
}

func TestAutoMigrateCreatesDefaultConfig(t *testing.T) {
	testIsolateEnv(t)
	resetGlobals()

	if err := AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}

	tomlPath := filepath.Join(dir, "config.toml")
	data, err := os.ReadFile(tomlPath)
	if err != nil {
		t.Fatalf("read toml: %v", err)
	}

	var got TOMLConfig
	md, err := toml.Decode(string(data), &got)
	if err != nil {
		t.Fatalf("decode toml: %v", err)
	}
	if undec := md.Undecoded(); len(undec) != 0 {
		t.Fatalf("undecoded fields: %v", undec)
	}

	expect := GetDefaultTOMLConfig()
	if diff := cmp.Diff(expect, got); diff != "" {
		t.Fatalf("default config mismatch (-want +got):\n%s", diff)
	}
}

func TestAutoMigrate(t *testing.T) {
	tests := []struct {
		name        string
		jsonContent []byte
		expect      TOMLConfig
	}{
		{
			name:        "standard configuration",
			jsonContent: standardJSON,
			expect: TOMLConfig{
				Base: BaseConfig{Emoji: true, FilesLength: 10},
				Alias: AliasConfig{
					Enabled: true,
					Entries: []AliasEntry{
						{Shortcuts: []string{"s", "feat", "feature"}, Name: "sparkles", Emoji: "✨"},
						{Shortcuts: []string{"c", "wip"}, Name: "construction", Emoji: "🚧"},
					},
				},
			},
		},
		{
			name:        "empty alias list",
			jsonContent: emptyAliasJSON,
			expect: TOMLConfig{
				Base:  BaseConfig{Emoji: true, FilesLength: 0},
				Alias: AliasConfig{Enabled: true, Entries: []AliasEntry{}},
			},
		},
		{
			name:        "missing fields",
			jsonContent: missingFieldsJSON,
			expect: TOMLConfig{
				Base: BaseConfig{Emoji: false, FilesLength: 0},
				Alias: AliasConfig{
					Enabled: false,
					Entries: []AliasEntry{{Shortcuts: []string{"b"}, Name: "bug", Emoji: "🐛"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testIsolateEnv(t)
			resetGlobals()

			dir, err := GetConfigDir()
			if err != nil {
				t.Fatalf("GetConfigDir: %v", err)
			}
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("MkdirAll: %v", err)
			}
			jsonPath := filepath.Join(dir, "config.json")
			if err := os.WriteFile(jsonPath, tt.jsonContent, 0644); err != nil {
				t.Fatalf("write json: %v", err)
			}

			if err := AutoMigrate(); err != nil {
				t.Fatalf("AutoMigrate: %v", err)
			}

			tomlPath := filepath.Join(dir, "config.toml")
			data, err := os.ReadFile(tomlPath)
			if err != nil {
				t.Fatalf("read toml: %v", err)
			}
			var got TOMLConfig
			if _, err := toml.Decode(string(data), &got); err != nil {
				t.Fatalf("decode toml: %v", err)
			}

			if diff := cmp.Diff(tt.expect.Base, got.Base); diff != "" {
				t.Fatalf("Base mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.expect.Alias, got.Alias); diff != "" {
				t.Fatalf("Alias mismatch (-want +got):\n%s", diff)
			}

			bakPath := filepath.Join(dir, "config.json.bak")
			bak, err := os.ReadFile(bakPath)
			if err != nil {
				t.Fatalf("backup not found: %v", err)
			}
			if diff := cmp.Diff(tt.jsonContent, bak); diff != "" {
				t.Fatalf("backup content mismatch (-want +got):\n%s", diff)
			}
			if _, err := os.Stat(jsonPath); err == nil {
				t.Fatalf("config.json should be renamed")
			}
		})
	}
}
