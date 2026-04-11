// NOTE: These tests mutate package-level globals.
// Do NOT use t.Parallel() in this file, as doing so may cause tests to interfere with
// each other and produce unreliable results. t.Parallel() is only safe if tests do not
// mutate shared state (such as package-level variables).
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/HidemaruOwO/pummit/internal/variable"
)

// isolateEnv prepares HOME for an isolated test and resets globals.
func isolateEnv(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, ".config"))

	prevCfg := CurrentTOMLConfig
	prevPath := TOMLConfigPath
	CurrentTOMLConfig = TOMLConfig{}
	TOMLConfigPath = ""

	t.Cleanup(func() {
		CurrentTOMLConfig = prevCfg
		TOMLConfigPath = prevPath
	})
}

// readConfigFile loads TOMLConfig from disk for verification.
func readConfigFile(t *testing.T) TOMLConfig {
	t.Helper()
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	p := filepath.Join(dir, "config.toml")
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("Stat: %v", err)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var cfg TOMLConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	return cfg
}

// TestTOMLConfigRoundTrip ensures saving and loading are idempotent.
func TestTOMLConfigRoundTrip(t *testing.T) {
	isolateEnv(t)
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	if err := SaveTOMLConfig(); err != nil {
		t.Fatalf("SaveTOMLConfig: %v", err)
	}
	CurrentTOMLConfig = TOMLConfig{}
	TOMLConfigPath = ""
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	want := GetDefaultTOMLConfig()
	if !reflect.DeepEqual(want, CurrentTOMLConfig) {
		t.Fatalf("CurrentTOMLConfig = %#v want %#v", CurrentTOMLConfig, want)
	}
	fileCfg := readConfigFile(t)
	if !reflect.DeepEqual(want, fileCfg) {
		t.Fatalf("fileCfg = %#v want %#v", fileCfg, want)
	}
}

// TestLoadTOMLConfigCreatesDefault verifies a default file is generated.
func TestLoadTOMLConfigCreatesDefault(t *testing.T) {
	isolateEnv(t)
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	want := GetDefaultTOMLConfig()
	if !reflect.DeepEqual(want, CurrentTOMLConfig) {
		t.Fatalf("CurrentTOMLConfig = %#v want %#v", CurrentTOMLConfig, want)
	}
	_ = readConfigFile(t)
}

// TestLoadTOMLConfigInvalid checks invalid TOML handling.
func TestLoadTOMLConfigInvalid(t *testing.T) {
	isolateEnv(t)
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte("::invalid::"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := LoadTOMLConfig(); err == nil {
		t.Fatalf("expected error for invalid TOML")
	} else {
		var parseErr toml.ParseError
		if !errors.As(err, &parseErr) {
			t.Fatalf("unexpected error type: %v", err)
		}
	}
}

// TestSaveTOMLConfigError covers directory creation failures.
func TestSaveTOMLConfigError(t *testing.T) {
	isolateEnv(t)
	home := os.Getenv("HOME")
	if err := os.WriteFile(filepath.Join(home, ".config"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	if err := SaveTOMLConfig(); err == nil {
		t.Fatalf("expected error when config dir is file")
	}
}

// TestLoadTOMLConfigCreateError covers failure when default file cannot be created.
func TestLoadTOMLConfigCreateError(t *testing.T) {
	isolateEnv(t)
	home := os.Getenv("HOME")
	if err := os.WriteFile(filepath.Join(home, ".config"), []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := LoadTOMLConfig(); err == nil {
		t.Fatalf("expected error when config dir is file")
	}
}

// TestLoadTOMLConfigUnknownKeys ensures unknown keys are ignored.
func TestLoadTOMLConfigUnknownKeys(t *testing.T) {
	isolateEnv(t)
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(GetDefaultTOMLConfig()); err != nil {
		t.Fatalf("Encode: %v", err)
	}
	buf.WriteString("unknown = 1\n")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	want := GetDefaultTOMLConfig()
	if !reflect.DeepEqual(want, CurrentTOMLConfig) {
		t.Fatalf("CurrentTOMLConfig = %#v want %#v", CurrentTOMLConfig, want)
	}
}

// TestLoadTOMLConfigPartial verifies missing fields are zero values.
func TestLoadTOMLConfigPartial(t *testing.T) {
	isolateEnv(t)
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	content := []byte("[base]\nemoji=false\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	// The design intentionally leaves unspecified fields at zero values
	// instead of merging with defaults.
	if CurrentTOMLConfig.Base.Emoji != false {
		t.Fatalf("Base.Emoji got %v", CurrentTOMLConfig.Base.Emoji)
	}
	if CurrentTOMLConfig.Base.FilesLength != 0 {
		t.Fatalf("Base.FilesLength want 0 got %d", CurrentTOMLConfig.Base.FilesLength)
	}
}

// TestSaveTOMLConfigPermissionError checks write permission failures.
func TestSaveTOMLConfigPermissionError(t *testing.T) {
	if !isUnixNonRoot() {
		t.Skip("skipping permission test: requires non-root Unix-like system")
	}
	isolateEnv(t)
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir: %v", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatalf("Chmod(file): %v", err)
	}
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatalf("Chmod(dir): %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(path, 0o644)
		_ = os.Chmod(dir, 0o755)
	})
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	if err := SaveTOMLConfig(); err == nil {
		t.Fatalf("expected permission error")
	}
}

// TestTOMLConfigPathOverride documents path override behavior.
func TestTOMLConfigPathOverride(t *testing.T) {
	isolateEnv(t)
	custom := filepath.Join(t.TempDir(), "custom.toml")
	TOMLConfigPath = custom
	CurrentTOMLConfig = GetDefaultTOMLConfig()
	if err := SaveTOMLConfig(); err != nil {
		t.Fatalf("SaveTOMLConfig: %v", err)
	}
	if _, err := os.Stat(custom); err != nil {
		t.Fatalf("Stat: %v", err)
	}
	CurrentTOMLConfig = TOMLConfig{}
	if err := LoadTOMLConfig(); err != nil {
		t.Fatalf("LoadTOMLConfig: %v", err)
	}
	// Specification A: explicit override path remains after Load.
	if TOMLConfigPath != custom {
		t.Fatalf("TOMLConfigPath = %s want %s", TOMLConfigPath, custom)
	}
}

func TestGetConfigDir(t *testing.T) {
	tests := []struct {
		name       string
		goos       string
		setup      func(t *testing.T)
		wantSuffix string
	}{
		{
			name: "unix default",
			setup: func(t *testing.T) {
				isolateEnv(t)
			},
			wantSuffix: filepath.Join(".config", "pummit"),
		},
		{
			name: "windows APPDATA",
			goos: "windows",
			setup: func(t *testing.T) {
				home := t.TempDir()
				t.Setenv("HOME", home)
				appData := filepath.Join(home, "AppData", "Roaming")
				if err := os.MkdirAll(appData, 0o755); err != nil {
					t.Fatalf("MkdirAll: %v", err)
				}
				t.Setenv("APPDATA", appData)
			},
			wantSuffix: filepath.Join("AppData", "Roaming", "pummit"),
		},
		{
			name: "windows fallback",
			goos: "windows",
			setup: func(t *testing.T) {
				home := t.TempDir()
				t.Setenv("HOME", home)
				t.Setenv("APPDATA", "")
			},
			wantSuffix: filepath.Join(".pummit"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.goos != "" && runtime.GOOS != tt.goos {
				t.Skipf("only runs on %s", tt.goos)
			}
			if tt.setup != nil {
				tt.setup(t)
			}
			got, err := GetConfigDir()
			if err != nil {
				t.Fatalf("GetConfigDir: %v", err)
			}
			if !strings.HasSuffix(got, tt.wantSuffix) {
				t.Fatalf("path %s does not end with %s", got, tt.wantSuffix)
			}
		})
	}
}

func TestInitCreatesDefaultTOML(t *testing.T) {
	isolateEnv(t)
	prevConfigPath := ConfigPath
	prevTOMLPath := TOMLConfigPath
	prevTOMLCfg := CurrentTOMLConfig
	t.Cleanup(func() {
		ConfigPath = prevConfigPath
		TOMLConfigPath = prevTOMLPath
		CurrentTOMLConfig = prevTOMLCfg
	})

	ConfigPath = ""
	TOMLConfigPath = ""
	CurrentTOMLConfig = TOMLConfig{}

	if err := Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}

	if ConfigPath == "" || TOMLConfigPath == "" {
		t.Fatal("ConfigPath and TOMLConfigPath should be set")
	}
	if !fileExists(TOMLConfigPath) {
		t.Fatalf("TOML config not created at %s", TOMLConfigPath)
	}
	if !strings.HasSuffix(ConfigPath, filepath.Join("pummit", "config.json")) {
		t.Fatalf("unexpected ConfigPath: %s", ConfigPath)
	}
	if !reflect.DeepEqual(CurrentTOMLConfig, GetDefaultTOMLConfig()) {
		t.Fatalf("CurrentTOMLConfig = %#v", CurrentTOMLConfig)
	}
}

func TestLoadAndSaveJSONConfig(t *testing.T) {
	isolateEnv(t)
	prevConfig := CurrentConfig
	prevDefault := DefaultConfig
	prevConfigPath := ConfigPath
	t.Cleanup(func() {
		CurrentConfig = prevConfig
		DefaultConfig = prevDefault
		ConfigPath = prevConfigPath
	})

	tempDir := t.TempDir()
	ConfigPath = filepath.Join(tempDir, "config.json")
	if err := os.MkdirAll(filepath.Dir(ConfigPath), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := Load(); err != nil {
		t.Fatalf("Load (create default): %v", err)
	}
	var expectedDefault Config
	if err := json.Unmarshal([]byte(variable.DEFAULT_CONFIG), &expectedDefault); err != nil {
		t.Fatalf("Unmarshal default: %v", err)
	}
	if !reflect.DeepEqual(CurrentConfig, expectedDefault) {
		t.Fatalf("CurrentConfig = %#v, want %#v", CurrentConfig, expectedDefault)
	}
	if !fileExists(ConfigPath) {
		t.Fatal("expected config.json to be created")
	}

	custom := Config{UseRawEmoji: true, UseAlias: true, UseFilesLength: true, FilesLength: 12, Aliases: [][]string{{"s", "sparkles", "✨"}}}
	CurrentConfig = custom
	if err := Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	CurrentConfig = Config{}
	if err := Load(); err != nil {
		t.Fatalf("Load (existing file): %v", err)
	}
	if !reflect.DeepEqual(CurrentConfig, custom) {
		t.Fatalf("CurrentConfig after reload = %#v, want %#v", CurrentConfig, custom)
	}
}

func TestConvertJSONToTOML(t *testing.T) {
	tests := []struct {
		name   string
		input  Config
		expect TOMLConfig
	}{
		{
			name: "full config",
			input: Config{
				UseRawEmoji:    true,
				UseAlias:       true,
				UseFilesLength: true,
				FilesLength:    50,
				Aliases:        [][]string{{"s,feat", "sparkles", "✨"}},
			},
			expect: TOMLConfig{
				Base: BaseConfig{Emoji: true, FilesLength: 50},
				Alias: AliasConfig{
					Enabled: true,
					Entries: []AliasEntry{{Shortcuts: []string{"s", "feat"}, Name: "sparkles", Emoji: "✨"}},
				},
			},
		},
		{
			name: "empty aliases",
			input: Config{
				UseRawEmoji: false,
				Aliases:     [][]string{},
			},
			expect: TOMLConfig{
				Base:  BaseConfig{Emoji: false, FilesLength: 0},
				Alias: AliasConfig{Enabled: false, Entries: []AliasEntry{}},
			},
		},
		{
			name: "no files length limit",
			input: Config{
				UseFilesLength: false,
				FilesLength:    100,
			},
			expect: TOMLConfig{
				Base:  BaseConfig{Emoji: false, FilesLength: 0},
				Alias: AliasConfig{Enabled: false, Entries: []AliasEntry{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertJSONToTOML(tt.input)
			if got.Meta.Version != "3.0" {
				t.Fatalf("Meta.Version = %s", got.Meta.Version)
			}
			if !reflect.DeepEqual(tt.expect.Base, got.Base) {
				t.Fatalf("Base = %#v, want %#v", got.Base, tt.expect.Base)
			}
			if tt.expect.Alias.Enabled != got.Alias.Enabled {
				t.Fatalf("Alias.Enabled = %v, want %v", got.Alias.Enabled, tt.expect.Alias.Enabled)
			}
			if !reflect.DeepEqual(tt.expect.Alias.Entries, got.Alias.Entries) {
				t.Fatalf("Alias entries = %#v, want %#v", got.Alias.Entries, tt.expect.Alias.Entries)
			}
		})
	}
}

func TestSplitCommaSeparated(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{"a, b, c", []string{"a", "b", "c"}},
		{"single", []string{"single"}},
		{"", []string{}},
		{"  a  ,  b  ", []string{"a", "b"}},
		{",,,", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitCommaSeparated(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTOMLToLegacyConfig(t *testing.T) {
	input := TOMLConfig{
		Base: BaseConfig{
			Emoji:       true,
			FilesLength: 50,
		},
		Alias: AliasConfig{
			Enabled: true,
			Entries: []AliasEntry{
				{Shortcuts: []string{"s", "feat"}, Name: "sparkles", Emoji: "✨"},
			},
		},
	}

	got := TOMLToLegacyConfig(input)

	if !got.UseRawEmoji {
		t.Error("UseRawEmoji should be true")
	}
	if !got.UseFilesLength {
		t.Error("UseFilesLength should be true")
	}
	if got.FilesLength != 50 {
		t.Errorf("FilesLength = %d, want 50", got.FilesLength)
	}
	if !got.UseAlias {
		t.Error("UseAlias should be true")
	}
	if len(got.Aliases) != 1 {
		t.Fatalf("Aliases length = %d, want 1", len(got.Aliases))
	}
	expectedAlias := []string{"s,feat", "sparkles", "✨"}
	if !reflect.DeepEqual(got.Aliases[0], expectedAlias) {
		t.Fatalf("Aliases[0] = %#v, want %#v", got.Aliases[0], expectedAlias)
	}
}

func TestGetAlias(t *testing.T) {
	t.Run("toml config", func(t *testing.T) {
		isolateEnv(t)
		dir, err := GetConfigDir()
		if err != nil {
			t.Fatalf("GetConfigDir: %v", err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		tomlPath := filepath.Join(dir, "config.toml")
		if err := os.WriteFile(tomlPath, []byte(""), 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		TOMLConfigPath = tomlPath
		CurrentTOMLConfig = TOMLConfig{
			Alias: AliasConfig{
				Entries: []AliasEntry{{Shortcuts: []string{"s", "feat"}, Name: "sparkles", Emoji: "✨"}},
			},
		}

		tests := []struct {
			name      string
			lookup    string
			wantName  string
			wantEmoji string
			wantFound bool
		}{
			{"found by shortcut", "s", "sparkles", "✨", true},
			{"found by name", "sparkles", "sparkles", "✨", true},
			{"not found", "unknown", "", "", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				name, emoji, found := GetAlias(tt.lookup)
				if found != tt.wantFound {
					t.Fatalf("found = %v, want %v", found, tt.wantFound)
				}
				if name != tt.wantName || emoji != tt.wantEmoji {
					t.Fatalf("got (%s, %s), want (%s, %s)", name, emoji, tt.wantName, tt.wantEmoji)
				}
			})
		}
	})

	t.Run("legacy json fallback", func(t *testing.T) {
		isolateEnv(t)
		prevConfig := CurrentConfig
		t.Cleanup(func() { CurrentConfig = prevConfig })

		CurrentConfig = Config{
			Aliases: [][]string{{"s,feat", "sparkles", "✨"}},
		}

		name, emoji, found := GetAlias("feat")
		if !found {
			t.Fatal("expected alias to be found")
		}
		if name != "sparkles" || emoji != "✨" {
			t.Fatalf("got (%s, %s)", name, emoji)
		}
	})
}

func TestAliasError(t *testing.T) {
	tests := []struct {
		errType string
		name    string
		want    string
	}{
		{"exists", "test", "alias 'test' already exists"},
		{"not_found", "test", "alias 'test' not found"},
		{"shortcut_exists", "s", "shortcut 's' already exists"},
		{"unknown", "x", "alias error: x"},
	}
	for _, tt := range tests {
		t.Run(tt.errType, func(t *testing.T) {
			err := &AliasError{Type: tt.errType, Name: tt.name}
			if err.Error() != tt.want {
				t.Fatalf("got %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

func TestMigrationError(t *testing.T) {
	tests := []struct {
		name string
		err  *MigrationError
		want string
	}{
		{
			name: "with underlying error",
			err:  &MigrationError{Op: "test", Path: "/path", Message: "failed", Err: errors.New("io error")},
			want: "migration test failed for /path: failed (io error)",
		},
		{
			name: "without underlying error",
			err:  &MigrationError{Op: "test", Path: "/path", Message: "failed"},
			want: "migration test failed for /path: failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Fatalf("got %q, want %q", tt.err.Error(), tt.want)
			}
		})
	}
}

func TestCheckConfigStatus(t *testing.T) {
	tests := []struct {
		name       string
		setupFiles func(dir string) error
		want       string
	}{
		{"none", func(dir string) error { return nil }, "none"},
		{"json_only", func(dir string) error {
			return os.WriteFile(filepath.Join(dir, "config.json"), []byte("{}"), 0o644)
		}, "json_only"},
		{"toml_only", func(dir string) error {
			return os.WriteFile(filepath.Join(dir, "config.toml"), []byte(""), 0o644)
		}, "toml_only"},
		{"both", func(dir string) error {
			if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte("{}"), 0o644); err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(dir, "config.toml"), []byte(""), 0o644)
		}, "both"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateEnv(t)
			dir, err := GetConfigDir()
			if err != nil {
				t.Fatalf("GetConfigDir: %v", err)
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatalf("MkdirAll: %v", err)
			}
			if err := tt.setupFiles(dir); err != nil {
				t.Fatalf("setupFiles: %v", err)
			}
			got, err := CheckConfigStatus()
			if err != nil {
				t.Fatalf("CheckConfigStatus: %v", err)
			}
			if got != tt.want {
				t.Fatalf("status = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMigrateConfig(t *testing.T) {
	t.Run("dry-run mode", func(t *testing.T) {
		isolateEnv(t)
		dir, err := GetConfigDir()
		if err != nil {
			t.Fatalf("GetConfigDir: %v", err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		jsonPath := filepath.Join(dir, "config.json")
		jsonContent := []byte(`{"writeEmoji":true,"useAlias":true,"useLimitPathesLength":true,"limitPathesLength":10,"alias":[["s,feat","sparkles","✨"]]}`)
		if err := os.WriteFile(jsonPath, jsonContent, 0o644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		result, err := MigrateConfig(false, true)
		if err != nil {
			t.Fatalf("MigrateConfig: %v", err)
		}
		if result.Success {
			t.Fatal("expected dry-run to be non-success")
		}
		if result.ConvertedFrom != "json" {
			t.Fatalf("ConvertedFrom = %s, want json", result.ConvertedFrom)
		}
		if result.BackupPath != "" {
			t.Fatalf("BackupPath = %s, want empty", result.BackupPath)
		}
		if !strings.Contains(result.Message, "[DRY-RUN]") {
			t.Fatalf("unexpected message: %s", result.Message)
		}
		if !fileExists(jsonPath) {
			t.Fatal("json file should remain in dry-run")
		}
		tomlPath := filepath.Join(dir, "config.toml")
		if fileExists(tomlPath) {
			t.Fatal("toml should not be created in dry-run")
		}
	})

	t.Run("force overwrites toml", func(t *testing.T) {
		isolateEnv(t)
		dir, err := GetConfigDir()
		if err != nil {
			t.Fatalf("GetConfigDir: %v", err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		jsonPath := filepath.Join(dir, "config.json")
		jsonContent := []byte(`{"writeEmoji":false,"useAlias":true,"useLimitPathesLength":true,"limitPathesLength":20,"alias":[["s,feat","sparkles","✨"]]}`)
		if err := os.WriteFile(jsonPath, jsonContent, 0o644); err != nil {
			t.Fatalf("WriteFile json: %v", err)
		}
		tomlPath := filepath.Join(dir, "config.toml")
		if err := os.WriteFile(tomlPath, []byte("base = { emoji = true }"), 0o644); err != nil {
			t.Fatalf("WriteFile toml: %v", err)
		}

		result, err := MigrateConfig(true, false)
		if err != nil {
			t.Fatalf("MigrateConfig: %v", err)
		}
		if !result.Success {
			t.Fatal("expected migration success")
		}
		if result.ConvertedFrom != "json" {
			t.Fatalf("ConvertedFrom = %s, want json", result.ConvertedFrom)
		}
		if result.BackupPath == "" || !fileExists(result.BackupPath) {
			t.Fatalf("backup not created: %s", result.BackupPath)
		}
		if fileExists(jsonPath) {
			t.Fatal("json file should be renamed after migration")
		}

		data, err := os.ReadFile(tomlPath)
		if err != nil {
			t.Fatalf("ReadFile toml: %v", err)
		}
		var got TOMLConfig
		if err := toml.Unmarshal(data, &got); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if got.Base.Emoji != false || got.Base.FilesLength != 20 {
			t.Fatalf("Base = %#v", got.Base)
		}
		if len(got.Alias.Entries) != 1 {
			t.Fatalf("expected one alias entry, got %d", len(got.Alias.Entries))
		}
		if got.Alias.Entries[0].Name != "sparkles" || got.Alias.Entries[0].Emoji != "✨" {
			t.Fatalf("alias entry = %#v", got.Alias.Entries[0])
		}
	})

	t.Run("creates default when json missing", func(t *testing.T) {
		isolateEnv(t)

		result, err := MigrateConfig(false, false)
		if err != nil {
			t.Fatalf("MigrateConfig: %v", err)
		}
		if !result.Success {
			t.Fatal("expected success when creating default")
		}
		if result.ConvertedFrom != "default" {
			t.Fatalf("ConvertedFrom = %s, want default", result.ConvertedFrom)
		}
		if result.Message != "Created default TOML config (no JSON config found)." {
			t.Fatalf("unexpected message: %s", result.Message)
		}

		fileCfg := readConfigFile(t)
		if !reflect.DeepEqual(fileCfg, GetDefaultTOMLConfig()) {
			t.Fatalf("default config mismatch: %#v", fileCfg)
		}
	})
}

func TestRollbackConfig(t *testing.T) {
	t.Run("successful rollback", func(t *testing.T) {
		isolateEnv(t)
		dir, err := GetConfigDir()
		if err != nil {
			t.Fatalf("GetConfigDir: %v", err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
		backupPath := filepath.Join(dir, "config.json.bak")
		backupContent := []byte(`{"writeEmoji":true}`)
		if err := os.WriteFile(backupPath, backupContent, 0o644); err != nil {
			t.Fatalf("WriteFile backup: %v", err)
		}
		tomlPath := filepath.Join(dir, "config.toml")
		if err := os.WriteFile(tomlPath, []byte("base = { emoji = true }"), 0o644); err != nil {
			t.Fatalf("WriteFile toml: %v", err)
		}

		if err := RollbackConfig(backupPath); err != nil {
			t.Fatalf("RollbackConfig: %v", err)
		}

		jsonPath := filepath.Join(dir, "config.json")
		data, err := os.ReadFile(jsonPath)
		if err != nil {
			t.Fatalf("ReadFile json: %v", err)
		}
		if string(data) != string(backupContent) {
			t.Fatalf("restored content mismatch: %s", string(data))
		}
		if fileExists(tomlPath) {
			t.Fatal("toml file should be removed on rollback")
		}
	})

	t.Run("missing backup", func(t *testing.T) {
		isolateEnv(t)
		err := RollbackConfig(filepath.Join(t.TempDir(), "missing.bak"))
		if err == nil {
			t.Fatal("expected error for missing backup")
		}
		var mErr *MigrationError
		if !errors.As(err, &mErr) {
			t.Fatalf("expected MigrationError, got %T", err)
		}
		if !strings.Contains(err.Error(), "backup file does not exist") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
