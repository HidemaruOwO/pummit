package alias

import (
	"errors"
	"fmt"

	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/internal/utils"
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

// findAlias searches for an alias by shortcut and returns its index and a
// pointer to the entry to enable in-place modifications.
func findAlias(shortcut string) (int, *config.AliasEntry) {
	// Iterate by index to avoid copying the slice value when taking its address.
	for i := 0; i < len(config.CurrentTOMLConfig.Alias.Entries); i++ {
		entry := &config.CurrentTOMLConfig.Alias.Entries[i]
		for _, s := range entry.Shortcuts {
			if s == shortcut {
				return i, entry
			}
		}
	}
	return -1, nil
}

// findEmojiIndex searches for an alias by emoji and returns its index.
func findEmojiIndex(emoji string) (int, bool) {
	for i, entry := range config.CurrentTOMLConfig.Alias.Entries {
		if entry.Emoji == emoji {
			return i, true
		}
	}
	return -1, false
}

// GetEmoji returns the prefix and emoji for a given alias name.
// Returns (found, prefix, emoji) where found indicates if the alias exists.
func GetEmoji(name string) (bool, string, string) {
	_, entry := findAlias(name)
	if entry == nil {
		return false, "", ""
	}
	return true, entry.Name, entry.Emoji
}

// Add adds a new alias with the given shortcut, prefix name, and emoji.
func Add(shortcut, name, emoji string) error {
	if _, existing := findAlias(shortcut); existing != nil {
		return ErrAliasExists
	}

	idx, exists := findEmojiIndex(emoji)

	if exists {
		// Emoji already exists, add the new shortcut to the existing entry.
		entry := &config.CurrentTOMLConfig.Alias.Entries[idx]
		if entry.Name != name {
			return fmt.Errorf("emoji '%s' already exists with a different name '%s', requested '%s'", emoji, entry.Name, name)
		}
		entry.Shortcuts = append(entry.Shortcuts, shortcut)
	} else {
		// Add a new entry for the new emoji.
		newEntry := config.AliasEntry{
			Shortcuts: []string{shortcut},
			Name:      name,
			Emoji:     emoji,
		}
		config.CurrentTOMLConfig.Alias.Entries = append(config.CurrentTOMLConfig.Alias.Entries, newEntry)
	}

	return config.SaveTOMLConfig()
}

// Delete removes an alias by shortcut.
func Delete(shortcut string) error {
	idx, entry := findAlias(shortcut)
	if entry == nil {
		return ErrAliasNotFound
	}

	// Remove the shortcut from the list.
	entry.Shortcuts = utils.RemoveString(entry.Shortcuts, shortcut)

	// If no shortcuts are left, remove the entire entry.
	if len(entry.Shortcuts) == 0 {
		config.CurrentTOMLConfig.Alias.Entries = append(config.CurrentTOMLConfig.Alias.Entries[:idx], config.CurrentTOMLConfig.Alias.Entries[idx+1:]...)
	}

	return config.SaveTOMLConfig()
}

// Get returns the emoji for a given alias name.
func Get(name string) (string, error) {
	_, entry := findAlias(name)
	if entry == nil {
		return "", ErrAliasNotFound
	}
	return entry.Emoji, nil
}

// List returns a map of all aliases (shortcuts) to their emoji values.
func List() map[string]string {
	result := make(map[string]string)
	for _, entry := range config.CurrentTOMLConfig.Alias.Entries {
		for _, shortcut := range entry.Shortcuts {
			result[shortcut] = entry.Emoji
		}
	}
	return result
}

// Reset resets the aliases to the default configuration.
func Reset() error {
	defaultTOML := config.GetDefaultTOMLConfig()
	config.CurrentTOMLConfig.Alias = defaultTOML.Alias
	return config.SaveTOMLConfig()
}
