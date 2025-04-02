package alias

import (
	"errors"
	"github.com/HidemaruOwO/pummit/internal/config"
	"strings"
)

type Alias struct {
	Name  string
	Emoji string
}

var (
	ErrAliasExists   = errors.New("alias already exists")
	ErrAliasNotFound = errors.New("alias not found")
)

func findAlias(name string) (int, bool) {
	for i, a := range config.CurrentConfig.Aliases {
		if a[0] == name {
			return i, true
		}

		if strings.Contains(a[0], ",") {
			parts := strings.Split(a[0], ",")
			for _, part := range parts {
				if part == name {
					return i, true
				}
			}
		}
	}
	return -1, false
}

func findEmojiIndex(emoji string) (int, bool) {
	for i, a := range config.CurrentConfig.Aliases {
		if len(a) >= 3 && a[2] == emoji {
			return i, true
		}
		if len(a) == 2 && a[1] == emoji {
			return i, true
		}
	}
	return -1, false
}

func GetEmoji(name string) (bool, string) {
	idx, exists := findAlias(name)
	if !exists {
		return false, ""
	}

	item := config.CurrentConfig.Aliases[idx]
	if len(item) >= 3 {
		return true, item[2]
	}
	return true, item[1]
}

func Add(name, emoji string) error {
	if _, exists := findAlias(name); exists {
		return ErrAliasExists
	}

	idx, exists := findEmojiIndex(emoji)

	if exists {
		item := config.CurrentConfig.Aliases[idx]
		parts := strings.Split(item[0], ",")

		for _, part := range parts {
			if part == name {
				return ErrAliasExists
			}
		}

		parts = append(parts, name)
		joined := strings.Join(parts, ",")

		if len(item) >= 3 {
			config.CurrentConfig.Aliases[idx] = []string{joined, item[1], item[2]}
		} else {
			config.CurrentConfig.Aliases[idx] = []string{joined, item[1]}
		}
	} else {
		config.CurrentConfig.Aliases = append(config.CurrentConfig.Aliases, []string{name, emoji})
	}

	return config.Save()
}

func Delete(name string) error {
	idx, exists := findAlias(name)
	if !exists {
		return ErrAliasNotFound
	}

	item := config.CurrentConfig.Aliases[idx]

	if item[0] == name {
		last := len(config.CurrentConfig.Aliases) - 1
		config.CurrentConfig.Aliases[idx] = config.CurrentConfig.Aliases[last]
		config.CurrentConfig.Aliases = config.CurrentConfig.Aliases[:last]
		return config.Save()
	}

	parts := strings.Split(item[0], ",")
	result := make([]string, 0, len(parts)-1)

	for _, part := range parts {
		if part != name {
			result = append(result, part)
		}
	}

	joined := strings.Join(result, ",")

	if len(item) >= 3 {
		config.CurrentConfig.Aliases[idx] = []string{joined, item[1], item[2]}
	} else {
		config.CurrentConfig.Aliases[idx] = []string{joined, item[1]}
	}

	return config.Save()
}

func Get(name string) (string, error) {
	idx, exists := findAlias(name)
	if !exists {
		return "", ErrAliasNotFound
	}

	item := config.CurrentConfig.Aliases[idx]
	if len(item) >= 3 {
		return item[2], nil
	}
	return item[1], nil
}

func List() map[string]string {
	result := make(map[string]string)
	for _, item := range config.CurrentConfig.Aliases {
		parts := strings.Split(item[0], ",")
		emoji := ""

		if len(item) >= 3 {
			emoji = item[2]
		} else {
			emoji = item[1]
		}

		for _, name := range parts {
			result[name] = emoji
		}
	}
	return result
}

func Reset() error {
	config.CurrentConfig.Aliases = make([][]string, 0)
	return config.Save()
}
