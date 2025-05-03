package alias

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/config"
)

// Alias represents an emoji alias with its name, prefix, and the actual emoji.
type Alias struct {
	Name   string
	Prefix string
	Emoji  string
}

var (
	// ErrAliasExists is returned when trying to add an alias that already exists.
	ErrAliasExists = errors.New("alias already exists")
	// ErrAliasNotFound is returned when an alias cannot be found.
	ErrAliasNotFound = errors.New("alias not found")
)

// findAlias searches for an alias by name and returns its index and existence.
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

// findEmojiIndex searches for an alias by emoji and returns its index and existence.
func findEmojiIndex(emoji string) (int, bool) {
	for i, a := range config.CurrentConfig.Aliases {
		if len(a) == 3 && a[2] == emoji {
			return i, true
		}
	}
	return -1, false
}

// GetEmoji returns the prefix and emoji for a given alias name.
// Returns (found, prefix, emoji) where found indicates if the alias exists.
func GetEmoji(name string) (bool, string, string) {
	idx, exists := findAlias(name)
	if !exists {
		return false, "", ""
	}

	item := config.CurrentConfig.Aliases[idx]
	if len(item) == 3 {
		return true, item[1], item[2]
	} else if len(item) == 2 {
		return true, "", item[1]
	}
	return false, "", ""
}

// Add adds a new alias with the given name, prefix, and emoji.
// If the alias already exists or the emoji is already registered with a different prefix,
// it returns an error.
func Add(name, prefix, emoji string) error {
	if _, exists := findAlias(name); exists {
		return ErrAliasExists
	}

	idx, exists := findEmojiIndex(emoji)

	if exists {
		item := config.CurrentConfig.Aliases[idx]
		if len(item) != 3 {
			return fmt.Errorf("inconsistent alias data format found for emoji '%s'", emoji)
		}

		if item[1] != prefix {
			return fmt.Errorf("emoji '%s' already exists with a different prefix '%s', requested '%s'", emoji, item[1], prefix)
		}

		parts := strings.Split(item[0], ",")

		for _, part := range parts {
			if part == name {
				return ErrAliasExists
			}
		}

		parts = append(parts, name)
		joined := strings.Join(parts, ",")

		config.CurrentConfig.Aliases[idx] = []string{joined, prefix, emoji}
	} else {
		config.CurrentConfig.Aliases = append(config.CurrentConfig.Aliases, []string{name, prefix, emoji})
	}

	return config.Save()
}

// Delete removes an alias by name.
// It returns ErrAliasNotFound if the alias doesn't exist.
func Delete(name string) error {
	idx, exists := findAlias(name)
	if !exists {
		return ErrAliasNotFound
	}

	item := config.CurrentConfig.Aliases[idx]

	if len(item) == 3 {
		if !strings.Contains(item[0], ",") {
			last := len(config.CurrentConfig.Aliases) - 1
			config.CurrentConfig.Aliases[idx] = config.CurrentConfig.Aliases[last]
			config.CurrentConfig.Aliases = config.CurrentConfig.Aliases[:last]
		} else {
			parts := strings.Split(item[0], ",")
			result := make([]string, 0, len(parts)-1)
			found := false
			for _, part := range parts {
				if part != name {
					result = append(result, part)
				} else {
					found = true
				}
			}
			if !found {
				return ErrAliasNotFound
			}
			joined := strings.Join(result, ",")
			config.CurrentConfig.Aliases[idx] = []string{joined, item[1], item[2]}
		}
	} else if len(item) == 2 {
		if !strings.Contains(item[0], ",") {
			last := len(config.CurrentConfig.Aliases) - 1
			config.CurrentConfig.Aliases[idx] = config.CurrentConfig.Aliases[last]
			config.CurrentConfig.Aliases = config.CurrentConfig.Aliases[:last]
		} else {
			parts := strings.Split(item[0], ",")
			result := make([]string, 0, len(parts)-1)
			found := false
			for _, part := range parts {
				if part != name {
					result = append(result, part)
				} else {
					found = true
				}
			}
			if !found {
				return ErrAliasNotFound
			}
			joined := strings.Join(result, ",")
			config.CurrentConfig.Aliases[idx] = []string{joined, item[1]}
		}
	} else {
		return fmt.Errorf("inconsistent alias data format found for alias '%s'", name)
	}

	return config.Save()
}

// Get returns the emoji for a given alias name.
// It returns an error if the alias doesn't exist.
func Get(name string) (string, error) {
	idx, exists := findAlias(name)
	if !exists {
		return "", ErrAliasNotFound
	}

	item := config.CurrentConfig.Aliases[idx]
	if len(item) == 3 {
		return item[2], nil
	} else if len(item) == 2 {
		return item[1], nil
	}
	return "", fmt.Errorf("inconsistent alias data format found for alias '%s'", name)
}

// List returns a map of all aliases to their emoji values.
func List() map[string]string {
	result := make(map[string]string)
	for _, item := range config.CurrentConfig.Aliases {
		aliasNames := ""
		emoji := ""

		if len(item) == 3 {
			aliasNames = item[0]
			emoji = item[2]
		} else if len(item) == 2 {
			aliasNames = item[0]
			emoji = item[1]
		} else {
			continue
		}

		parts := strings.Split(aliasNames, ",")
		for _, name := range parts {
			result[name] = emoji
		}
	}
	return result
}

// Reset resets the aliases to the default configuration.
func Reset() error {
	config.CurrentConfig.Aliases = config.DefaultConfig.Aliases
	return config.Save()
}
