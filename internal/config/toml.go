package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// TOML設定ファイルの構造体
type TOMLConfig struct {
	Meta        MetaConfig        `toml:"meta"`
	Base        BaseConfig        `toml:"base"`
	Interactive InteractiveConfig `toml:"interactive"`
	Locale      LocaleConfig      `toml:"locale"`
	Templates   TemplatesConfig   `toml:"templates"`
	Scope       ScopeConfig       `toml:"scope"`
	Alias       AliasConfig       `toml:"alias"`
	Branch      BranchConfig      `toml:"branchMapping"`
}

type MetaConfig struct {
	Version string `toml:"version"`
}

type BaseConfig struct {
	Emoji       bool `toml:"emoji"`
	FilesLength int  `toml:"filesLength"`
}

type InteractiveConfig struct {
	Enabled     bool   `toml:"enabled"`
	DefaultMode string `toml:"defaultMode"`
	ShowPreview bool   `toml:"showPreview"`
	FuzzySearch bool   `toml:"fuzzySearch"`
}

type LocaleConfig struct {
	Language   string `toml:"language"`
	AutoDetect bool   `toml:"autoDetect"`
}

type TemplatesConfig struct {
	Enabled         bool                          `toml:"enabled"`
	DefaultTemplate string                        `toml:"defaultTemplate"`
	Definitions     map[string]TemplateDefinition `toml:"definitions"`
}

type TemplateDefinition struct {
	Format             string `toml:"format"`
	Description        string `toml:"description"`
	DefaultEmoji       string `toml:"defaultEmoji,omitempty"`
	Scope              bool   `toml:"scope,omitempty"`
	RequireDescription bool   `toml:"requireDescription,omitempty"`
}

type ScopeConfig struct {
	Enabled      bool     `toml:"enabled"`
	AutoDetect   bool     `toml:"autoDetect"`
	Suggestions  []string `toml:"suggestions"`
	FromHistory  bool     `toml:"fromHistory"`
	HistoryLimit int      `toml:"historyLimit"`
}

type AliasConfig struct {
	Enabled bool         `toml:"enabled"`
	Entries []AliasEntry `toml:"entries"`
}

type AliasEntry struct {
	Shortcuts []string `toml:"shortcuts"`
	Name      string   `toml:"name"`
	Emoji     string   `toml:"emoji"`
}

type BranchConfig struct {
	Enabled  bool         `toml:"enabled"`
	Fallback string       `toml:"fallback"`
	Rules    []BranchRule `toml:"rules"`
}

type BranchRule struct {
	Pattern     string `toml:"pattern"`
	Emoji       string `toml:"emoji"`
	Description string `toml:"description"`
}

var (
	TOMLConfigPath    string
	CurrentTOMLConfig TOMLConfig
)

// デフォルトTOML設定を生成
func GetDefaultTOMLConfig() TOMLConfig {
	return TOMLConfig{
		Meta: MetaConfig{
			Version: "3.0",
		},
		Base: BaseConfig{
			Emoji:       true,
			FilesLength: 50,
		},
		Interactive: InteractiveConfig{
			Enabled:     true,
			DefaultMode: "emoji",
			ShowPreview: true,
			FuzzySearch: true,
		},
		Locale: LocaleConfig{
			Language:   "ja",
			AutoDetect: true,
		},
		Templates: TemplatesConfig{
			Enabled:         true,
			DefaultTemplate: "default",
			Definitions: map[string]TemplateDefinition{
				"default": {
					Format:      "{emoji} {message} ({files})",
					Description: "Default template",
				},
				"feat": {
					Format:             "{emoji} {scope}: {message}\n\n{description}\n\n({files})",
					Description:        "Feature addition template",
					DefaultEmoji:       "sparkles",
					Scope:              true,
					RequireDescription: true,
				},
			},
		},
		Scope: ScopeConfig{
			Enabled:      true,
			AutoDetect:   true,
			Suggestions:  []string{"api", "ui", "core", "auth", "db"},
			FromHistory:  true,
			HistoryLimit: 50,
		},
		Alias: AliasConfig{
			Enabled: true,
			Entries: []AliasEntry{
				{Shortcuts: []string{"s", "feat", "feature"}, Name: "sparkles", Emoji: "✨"},
				{Shortcuts: []string{"c", "wip"}, Name: "construction", Emoji: "🚧"},
				{Shortcuts: []string{"t", "new", "init"}, Name: "tada", Emoji: "🎉"},
				{Shortcuts: []string{"r", "pr", "pull", "merge"}, Name: "recycle", Emoji: "♻️"},
				{Shortcuts: []string{"wb", "rm", "remove", "del", "delete"}, Name: "wastebasket", Emoji: "🗑️"},
				{Shortcuts: []string{"b", "fix", "error"}, Name: "bug", Emoji: "🐛"},
				{Shortcuts: []string{"e", "lint", "format", "refactor"}, Name: "eyes", Emoji: "👀"},
				{Shortcuts: []string{"d", "doc", "docs", "document", "documents"}, Name: "books", Emoji: "📚"},
				{Shortcuts: []string{"a", "ui", "design", "icon", "icons"}, Name: "art", Emoji: "🎨"},
				{Shortcuts: []string{"h", "tune", "tuning", "perform", "performance"}, Name: "horse", Emoji: "🐎"},
				{Shortcuts: []string{"w", "change", "tool", "tools", "lib", "library"}, Name: "wrench", Emoji: "🔧"},
				{Shortcuts: []string{"l", "test", "testing"}, Name: "rotating_light", Emoji: "🚨"},
				{Shortcuts: []string{"sm", "special", "important"}, Name: "snowman", Emoji: "☃️"},
				{Shortcuts: []string{"p", "pack", "mod", "module"}, Name: "package", Emoji: "📦️"},
			},
		},
		Branch: BranchConfig{
			Enabled:  true,
			Fallback: "construction",
			Rules: []BranchRule{
				{Pattern: "^feature/.*", Emoji: "sparkles", Description: "Feature branch"},
				{Pattern: "^(fix|bugfix|hotfix)/.*", Emoji: "bug", Description: "Bug fix branch"},
				{Pattern: "^docs/.*", Emoji: "books", Description: "Documentation branch"},
				{Pattern: "^refactor/.*", Emoji: "eyes", Description: "Refactoring branch"},
				{Pattern: "^test/.*", Emoji: "rotating_light", Description: "Test branch"},
			},
		},
	}
}

// TOML設定ファイルを読み込み
func LoadTOMLConfig() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	TOMLConfigPath = filepath.Join(configDir, "config.toml")

	if _, err := os.Stat(TOMLConfigPath); os.IsNotExist(err) {
		// ファイルが存在しない場合はデフォルト設定で作成
		CurrentTOMLConfig = GetDefaultTOMLConfig()
		return SaveTOMLConfig()
	}

	data, err := os.ReadFile(TOMLConfigPath)
	if err != nil {
		return err
	}

	return toml.Unmarshal(data, &CurrentTOMLConfig)
}

// TOML設定ファイルを保存
func SaveTOMLConfig() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// 設定ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// TOMLConfigPathが設定されていない場合は設定
	if TOMLConfigPath == "" {
		TOMLConfigPath = filepath.Join(configDir, "config.toml")
	}

	file, err := os.Create(TOMLConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(CurrentTOMLConfig)
}

// レガシーJSONからTOML設定に変換
func ConvertJSONToTOML(jsonConfig Config) TOMLConfig {
	tomlConfig := GetDefaultTOMLConfig()

	// 基本設定の変換
	tomlConfig.Base.Emoji = jsonConfig.UseRawEmoji
	tomlConfig.Alias.Enabled = jsonConfig.UseAlias

	// ファイル長の制限設定
	if jsonConfig.UseFilesLength {
		tomlConfig.Base.FilesLength = jsonConfig.FilesLength
	} else {
		tomlConfig.Base.FilesLength = 0 // 制限なし
	}

	// エイリアスの変換
	tomlConfig.Alias.Entries = []AliasEntry{}
	for _, alias := range jsonConfig.Aliases {
		if len(alias) >= 3 {
			// JSON形式: ["s,feat,feature", "sparkles", "✨"]
			shortcuts := alias[0]
			name := alias[1]
			emoji := alias[2]

			// ショートカットをカンマで分割
			shortcutList := []string{}
			for _, shortcut := range splitCommaSeparated(shortcuts) {
				if shortcut != "" {
					shortcutList = append(shortcutList, shortcut)
				}
			}

			entry := AliasEntry{
				Shortcuts: shortcutList,
				Name:      name,
				Emoji:     emoji,
			}
			tomlConfig.Alias.Entries = append(tomlConfig.Alias.Entries, entry)
		}
	}

	// バージョンを設定
	tomlConfig.Meta.Version = "3.0"

	return tomlConfig
}

// カンマ区切り文字列を分割するヘルパー関数
func splitCommaSeparated(s string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	for _, part := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
