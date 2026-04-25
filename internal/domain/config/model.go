package config

type Config struct {
	Meta          MetaConfig        `toml:"meta"`
	Base          BaseConfig        `toml:"base"`
	Interactive   InteractiveConfig `toml:"interactive"`
	Locale        LocaleConfig      `toml:"locale"`
	Templates     TemplatesConfig   `toml:"templates"`
	Scope         ScopeConfig       `toml:"scope"`
	Alias         AliasConfig       `toml:"alias"`
	BranchMapping BranchConfig      `toml:"branchMapping"`
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
