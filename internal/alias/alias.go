package alias

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/internal/emojis"
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

func findEmojiIndex(actualEmoji string) (int, bool) {
	for i, a := range config.CurrentConfig.Aliases {
		if len(a) == 3 && a[2] == actualEmoji {
			return i, true
		}
	}
	return -1, false
}

// (e.g) found, prefix, emoji := alias.GetEmoji(cm.Emoji)
func GetEmoji(name string) (bool, string, string) {
	idx, exists := findAlias(name)
	if !exists {
		return false, "", ""
	}

	item := config.CurrentConfig.Aliases[idx]
	return true, item[1], item[2]
}

func Add(name, emojiName string) error {
	if _, exists := findAlias(name); exists {
		return ErrAliasExists
	}

	actualEmoji, err := emojis.GetEmojiByName(emojiName)
	if err != nil {
		return fmt.Errorf("failed to find emoji with name '%s': %w", emojiName, err)
	}

	idx, exists := findEmojiIndex(actualEmoji)

	if exists {
		item := config.CurrentConfig.Aliases[idx]

		if len(item) != 3 {
			return fmt.Errorf("inconsistent alias data format found for emoji '%s'", actualEmoji)
		}

		if item[1] != emojiName {
			return fmt.Errorf("emoji '%s' already exists with a different prefix '%s', requested '%s'", actualEmoji, item[1], emojiName)
		}

		parts := strings.Split(item[0], ",")

		for _, part := range parts {
			if part == name {
				return ErrAliasExists
			}
		}

		parts = append(parts, name)
		joined := strings.Join(parts, ",")

		config.CurrentConfig.Aliases[idx] = []string{joined, emojiName, actualEmoji}

	} else {
		config.CurrentConfig.Aliases = append(config.CurrentConfig.Aliases, []string{name, emojiName, actualEmoji})
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
	config.CurrentConfig.Aliases = config.DefaultConfig.Aliases
	return config.Save()
}
