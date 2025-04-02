package alias

import (
	"errors"
	"github.com/HidemaruOwO/pummit/internal/config"
)

type Alias struct {
	Name  string
	Emoji string
}

var (
	ErrAliasExists   = errors.New("alias already exists")
	ErrAliasNotFound = errors.New("alias not found")
)

// 指定されたエイリアスのインデックスを返す
func findAlias(name string) (int, bool) {
	for i, a := range config.CurrentConfig.Aliases {
		if a[0] == name {
			return i, true
		}
	}
	return -1, false
}

func Add(name, emoji string) error {
	if _, exists := findAlias(name); exists {
		return ErrAliasExists
	}
	config.CurrentConfig.Aliases = append(config.CurrentConfig.Aliases, []string{name, emoji})
	return config.Save()
}

func Delete(name string) error {
	idx, exists := findAlias(name)
	if !exists {
		return ErrAliasNotFound
	}

	last := len(config.CurrentConfig.Aliases) - 1
	config.CurrentConfig.Aliases[idx] = config.CurrentConfig.Aliases[last]
	config.CurrentConfig.Aliases = config.CurrentConfig.Aliases[:last]

	return config.Save()
}

func Get(name string) (string, error) {
	if idx, exists := findAlias(name); exists {
		return config.CurrentConfig.Aliases[idx][1], nil
	}
	return "", ErrAliasNotFound
}

func List() map[string]string {
	aliases := make(map[string]string, len(config.CurrentConfig.Aliases))
	for _, a := range config.CurrentConfig.Aliases {
		aliases[a[0]] = a[1]
	}
	return aliases
}

func Reset() error {
	config.CurrentConfig.Aliases = make([][]string, 0)
	return config.Save()
}
