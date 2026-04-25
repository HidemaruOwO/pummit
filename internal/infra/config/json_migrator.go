package config

import (
	"encoding/json"
	"strings"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
	"github.com/HidemaruOwO/pummit/internal/variable"
)

type legacyConfig struct {
	UseRawEmoji    bool       `json:"writeEmoji"`
	UseAlias       bool       `json:"useAlias"`
	UseFilesLength bool       `json:"useLimitPathesLength"`
	FilesLength    int        `json:"limitPathesLength"`
	Aliases        [][]string `json:"alias"`
}

func MigrateLegacyJSON(data []byte) (domainconfig.Config, error) {
	legacy, err := decodeLegacyJSON(data)
	if err != nil {
		return domainconfig.Config{}, err
	}

	cfg := domainconfig.Default()
	cfg.Base.Emoji = legacy.UseRawEmoji
	if legacy.UseFilesLength {
		cfg.Base.FilesLength = legacy.FilesLength
	} else {
		cfg.Base.FilesLength = 0
	}
	cfg.Alias.Enabled = legacy.UseAlias
	cfg.Alias.Entries = make([]domainconfig.AliasEntry, 0, len(legacy.Aliases))

	for _, raw := range legacy.Aliases {
		if len(raw) < 3 {
			continue
		}

		entry := domainconfig.AliasEntry{
			Shortcuts: splitShortcuts(raw[0]),
			Name:      raw[1],
			Emoji:     raw[2],
		}
		cfg.Alias.Entries = append(cfg.Alias.Entries, entry)
	}

	if err := domainconfig.Validate(cfg); err != nil {
		return domainconfig.Config{}, err
	}

	return cfg, nil
}

func decodeLegacyJSON(data []byte) (legacyConfig, error) {
	var cfg legacyConfig
	if err := json.Unmarshal([]byte(variable.DEFAULT_CONFIG), &cfg); err != nil {
		return legacyConfig{}, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return legacyConfig{}, err
	}

	return cfg, nil
}

func splitShortcuts(raw string) []string {
	parts := strings.Split(raw, ",")
	shortcuts := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		shortcuts = append(shortcuts, part)
	}

	return shortcuts
}
