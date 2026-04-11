package config

import "strings"

// TOML設定からレガシーJSON設定への変換（後方互換性のため）
func TOMLToLegacyConfig(tomlConfig TOMLConfig) Config {
	legacyConfig := Config{
		UseRawEmoji:    tomlConfig.Base.Emoji,
		UseAlias:       tomlConfig.Alias.Enabled,
		UseFilesLength: tomlConfig.Base.FilesLength > 0,
		FilesLength:    tomlConfig.Base.FilesLength,
		Aliases:        [][]string{},
	}

	// TOMLエイリアスをレガシー形式に変換
	for _, entry := range tomlConfig.Alias.Entries {
		if len(entry.Shortcuts) > 0 && entry.Name != "" && entry.Emoji != "" {
			// ショートカットをカンマ区切りで結合
			shortcuts := strings.Join(entry.Shortcuts, ",")
			alias := []string{shortcuts, entry.Name, entry.Emoji}
			legacyConfig.Aliases = append(legacyConfig.Aliases, alias)
		}
	}

	return legacyConfig
}

// レガシーJSON設定の互換性確保のためのラッパー関数群

// GetAlias はエイリアスを取得する（統合版）
func GetAlias(name string) (string, string, bool) {
	// TOML設定が有効な場合
	if TOMLConfigPath != "" && fileExists(TOMLConfigPath) {
		for _, entry := range CurrentTOMLConfig.Alias.Entries {
			for _, shortcut := range entry.Shortcuts {
				if shortcut == name {
					return entry.Name, entry.Emoji, true
				}
			}
			if entry.Name == name {
				return entry.Name, entry.Emoji, true
			}
		}
		return "", "", false
	}

	// フォールバック: レガシーJSON形式
	for _, alias := range CurrentConfig.Aliases {
		if len(alias) >= 3 {
			shortcuts := alias[0]
			prefix := alias[1]
			emoji := alias[2]

			// ショートカットをチェック
			for _, shortcut := range strings.Split(shortcuts, ",") {
				if strings.TrimSpace(shortcut) == name {
					return prefix, emoji, true
				}
			}

			// プレフィックス名をチェック
			if prefix == name {
				return prefix, emoji, true
			}
		}
	}

	return "", "", false
}

// エイリアス関連のエラー
type AliasError struct {
	Type string
	Name string
}

func (e *AliasError) Error() string {
	switch e.Type {
	case "exists":
		return "alias '" + e.Name + "' already exists"
	case "not_found":
		return "alias '" + e.Name + "' not found"
	case "shortcut_exists":
		return "shortcut '" + e.Name + "' already exists"
	default:
		return "alias error: " + e.Name
	}
}
