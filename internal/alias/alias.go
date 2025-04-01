package alias

import (
	"errors"

	"github.com/HidemaruOwO/pummit/internal/config"
)

// エラー定義
var (
	ErrAliasExists   = errors.New("alias already exists")
	ErrAliasNotFound = errors.New("alias not found")
)

// Add はエイリアスを追加します
func Add(name, emoji string) error {
	// 既に存在するか確認
	if _, exists := config.CurrentConfig.Aliases[name]; exists {
		return ErrAliasExists
	}

	// エイリアスを追加して保存
	config.CurrentConfig.Aliases[name] = emoji
	return config.Save()
}

// Delete はエイリアスを削除します
func Delete(name string) error {
	// エイリアスが存在するか確認
	if _, exists := config.CurrentConfig.Aliases[name]; !exists {
		return ErrAliasNotFound
	}

	// 削除して保存
	delete(config.CurrentConfig.Aliases, name)
	return config.Save()
}

// Get はエイリアスに関連付けられた絵文字を取得します
func Get(name string) (string, error) {
	emoji, exists := config.CurrentConfig.Aliases[name]
	if !exists {
		return "", ErrAliasNotFound
	}
	return emoji, nil
}

// List は全てのエイリアスのリストを返します
func List() map[string]string {
	return config.CurrentConfig.Aliases
}

// Reset は全てのエイリアスをリセットします
func Reset() error {
	config.CurrentConfig.Aliases = make(map[string]string)
	return config.Save()
}
