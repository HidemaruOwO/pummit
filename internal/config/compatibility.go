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

// GetCurrentConfig は現在の設定を取得する（後方互換性）
func GetCurrentConfig() Config {
	// TOML設定が利用可能な場合はそれを変換
	if TOMLConfigPath != "" && fileExists(TOMLConfigPath) {
		return TOMLToLegacyConfig(CurrentTOMLConfig)
	}
	// フォールバック: レガシーJSON設定を返す
	return CurrentConfig
}

// UpdateCurrentConfig は現在の設定を更新する（後方互換性）
func UpdateCurrentConfig(newConfig Config) error {
	// TOML設定が有効な場合はTOML形式で保存
	if TOMLConfigPath != "" && fileExists(TOMLConfigPath) {
		CurrentTOMLConfig = ConvertJSONToTOML(newConfig)
		return SaveTOMLConfig()
	}

	// フォールバック: レガシーJSON形式で保存
	CurrentConfig = newConfig
	return Save()
}

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

// AddAlias はエイリアスを追加する（統合版）
func AddAlias(shortcuts []string, name string, emoji string) error {
	// TOML設定が有効な場合
	if TOMLConfigPath != "" && fileExists(TOMLConfigPath) {
		// 既存のエイリアスをチェック
		for _, entry := range CurrentTOMLConfig.Alias.Entries {
			if entry.Name == name {
				return &AliasError{Type: "exists", Name: name}
			}
			for _, shortcut := range shortcuts {
				for _, existingShortcut := range entry.Shortcuts {
					if existingShortcut == shortcut {
						return &AliasError{Type: "shortcut_exists", Name: shortcut}
					}
				}
			}
		}

		// 新しいエイリアスを追加
		newEntry := AliasEntry{
			Shortcuts: shortcuts,
			Name:      name,
			Emoji:     emoji,
		}
		CurrentTOMLConfig.Alias.Entries = append(CurrentTOMLConfig.Alias.Entries, newEntry)
		return SaveTOMLConfig()
	}

	// フォールバック: レガシーJSON形式
	shortcutsStr := strings.Join(shortcuts, ",")

	// 既存のエイリアスをチェック
	for _, alias := range CurrentConfig.Aliases {
		if len(alias) >= 2 && alias[1] == name {
			return &AliasError{Type: "exists", Name: name}
		}
	}

	// 新しいエイリアスを追加
	newAlias := []string{shortcutsStr, name, emoji}
	CurrentConfig.Aliases = append(CurrentConfig.Aliases, newAlias)
	return Save()
}

// RemoveAlias はエイリアスを削除する（統合版）
func RemoveAlias(name string) error {
	// TOML設定が有効な場合
	if TOMLConfigPath != "" && fileExists(TOMLConfigPath) {
		found := false
		newEntries := []AliasEntry{}

		for _, entry := range CurrentTOMLConfig.Alias.Entries {
			if entry.Name != name {
				newEntries = append(newEntries, entry)
			} else {
				found = true
			}
		}

		if !found {
			return &AliasError{Type: "not_found", Name: name}
		}

		CurrentTOMLConfig.Alias.Entries = newEntries
		return SaveTOMLConfig()
	}

	// フォールバック: レガシーJSON形式
	found := false
	newAliases := [][]string{}

	for _, alias := range CurrentConfig.Aliases {
		if len(alias) >= 2 && alias[1] != name {
			newAliases = append(newAliases, alias)
		} else {
			found = true
		}
	}

	if !found {
		return &AliasError{Type: "not_found", Name: name}
	}

	CurrentConfig.Aliases = newAliases
	return Save()
}

// ListAliases はエイリアス一覧を取得する（統合版）
func ListAliases() []AliasEntry {
	// TOML設定が有効な場合
	if TOMLConfigPath != "" && fileExists(TOMLConfigPath) {
		return CurrentTOMLConfig.Alias.Entries
	}

	// フォールバック: レガシーJSON形式から変換
	var entries []AliasEntry
	for _, alias := range CurrentConfig.Aliases {
		if len(alias) >= 3 {
			shortcuts := strings.Split(alias[0], ",")
			for i, shortcut := range shortcuts {
				shortcuts[i] = strings.TrimSpace(shortcut)
			}

			entry := AliasEntry{
				Shortcuts: shortcuts,
				Name:      alias[1],
				Emoji:     alias[2],
			}
			entries = append(entries, entry)
		}
	}
	return entries
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
