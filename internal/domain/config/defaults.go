package config

func Default() Config {
	return Config{
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
				{Shortcuts: []string{"h", "tune", "tuning", "perf", "perform", "performance"}, Name: "rocket", Emoji: "🚀"},
				{Shortcuts: []string{"w", "change", "tool", "tools", "lib", "library"}, Name: "wrench", Emoji: "🔧"},
				{Shortcuts: []string{"l", "test", "testing"}, Name: "rotating_light", Emoji: "🚨"},
				{Shortcuts: []string{"p", "pack", "mod", "module"}, Name: "package", Emoji: "📦️"},
			},
		},
		BranchMapping: BranchConfig{
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
