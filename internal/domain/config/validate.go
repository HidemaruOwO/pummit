package config

import (
	"fmt"
	"regexp"
)

var validModes = map[string]struct{}{
	"emoji":    {},
	"template": {},
	"scope":    {},
	"branch":   {},
}

var validLanguages = map[string]struct{}{
	"ja": {},
	"en": {},
}

func Validate(cfg Config) error {
	if cfg.Meta.Version == "" {
		return fmt.Errorf("meta.version is required")
	}

	if cfg.Base.FilesLength < 0 {
		return fmt.Errorf("base.filesLength must be >= 0")
	}

	if _, ok := validModes[cfg.Interactive.DefaultMode]; !ok {
		return fmt.Errorf("interactive.defaultMode must be one of emoji, template, scope, branch")
	}

	if _, ok := validLanguages[cfg.Locale.Language]; !ok {
		return fmt.Errorf("locale.language must be one of ja, en")
	}

	if cfg.Scope.HistoryLimit < 0 {
		return fmt.Errorf("scope.historyLimit must be >= 0")
	}

	if cfg.Templates.DefaultTemplate == "" {
		return fmt.Errorf("templates.defaultTemplate is required")
	}

	if len(cfg.Templates.Definitions) == 0 {
		return fmt.Errorf("templates.definitions must not be empty")
	}

	if cfg.Templates.Enabled {
		if _, ok := cfg.Templates.Definitions[cfg.Templates.DefaultTemplate]; !ok {
			return fmt.Errorf("templates.defaultTemplate %q is not defined", cfg.Templates.DefaultTemplate)
		}
	}

	for name, def := range cfg.Templates.Definitions {
		if def.Format == "" {
			return fmt.Errorf("templates.definitions.%s.format is required", name)
		}
		if def.Description == "" {
			return fmt.Errorf("templates.definitions.%s.description is required", name)
		}
	}

	seenShortcuts := map[string]struct{}{}
	for i, entry := range cfg.Alias.Entries {
		if entry.Name == "" {
			return fmt.Errorf("alias.entries[%d].name is required", i)
		}
		if entry.Emoji == "" {
			return fmt.Errorf("alias.entries[%d].emoji is required", i)
		}
		if len(entry.Shortcuts) == 0 {
			return fmt.Errorf("alias.entries[%d].shortcuts must not be empty", i)
		}
		for _, shortcut := range entry.Shortcuts {
			if shortcut == "" {
				return fmt.Errorf("alias.entries[%d].shortcuts must not contain empty values", i)
			}
			if _, ok := seenShortcuts[shortcut]; ok {
				return fmt.Errorf("alias shortcut %q is duplicated", shortcut)
			}
			seenShortcuts[shortcut] = struct{}{}
		}
	}

	if cfg.BranchMapping.Fallback == "" {
		return fmt.Errorf("branchMapping.fallback is required")
	}

	for i, rule := range cfg.BranchMapping.Rules {
		if rule.Pattern == "" {
			return fmt.Errorf("branchMapping.rules[%d].pattern is required", i)
		}
		if _, err := regexp.Compile(rule.Pattern); err != nil {
			return fmt.Errorf("branchMapping.rules[%d].pattern is invalid: %w", i, err)
		}
		if rule.Emoji == "" {
			return fmt.Errorf("branchMapping.rules[%d].emoji is required", i)
		}
		if rule.Description == "" {
			return fmt.Errorf("branchMapping.rules[%d].description is required", i)
		}
	}

	return nil
}
