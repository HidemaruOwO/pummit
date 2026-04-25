package usecase

import (
	"fmt"
	"sort"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
)

type AliasService struct {
	config  *ConfigService
	catalog Catalog
	remote  GitmojiLookup
}

func NewAliasService(config *ConfigService, catalog Catalog, remote GitmojiLookup) *AliasService {
	return &AliasService{config: config, catalog: catalog, remote: remote}
}

func (s *AliasService) Add(shortcut, name, explicitEmoji string) error {
	cfg, err := s.config.Load()
	if err != nil {
		return err
	}

	for _, entry := range cfg.Alias.Entries {
		for _, existing := range entry.Shortcuts {
			if existing == shortcut {
				return fmt.Errorf("alias shortcut %q already exists", shortcut)
			}
		}
	}

	emojiValue, err := s.resolveEmoji(name, explicitEmoji)
	if err != nil {
		return err
	}

	matched := false
	for i, entry := range cfg.Alias.Entries {
		if entry.Name == name || entry.Emoji == emojiValue {
			if entry.Name != name {
				return fmt.Errorf("emoji %q already belongs to alias %q", emojiValue, entry.Name)
			}
			cfg.Alias.Entries[i].Shortcuts = append(cfg.Alias.Entries[i].Shortcuts, shortcut)
			matched = true
			break
		}
	}

	if !matched {
		cfg.Alias.Entries = append(cfg.Alias.Entries, domainconfig.AliasEntry{
			Shortcuts: []string{shortcut},
			Name:      name,
			Emoji:     emojiValue,
		})
	}

	if err := domainconfig.Validate(cfg); err != nil {
		return err
	}

	return s.config.Save(cfg)
}

func (s *AliasService) List() ([]domainconfig.AliasEntry, error) {
	cfg, err := s.config.Load()
	if err != nil {
		return nil, err
	}

	entries := append([]domainconfig.AliasEntry(nil), cfg.Alias.Entries...)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	return entries, nil
}

func (s *AliasService) Delete(shortcut string) error {
	cfg, err := s.config.Load()
	if err != nil {
		return err
	}

	found := false
	updated := make([]domainconfig.AliasEntry, 0, len(cfg.Alias.Entries))
	for _, entry := range cfg.Alias.Entries {
		shortcuts := make([]string, 0, len(entry.Shortcuts))
		for _, existing := range entry.Shortcuts {
			if existing == shortcut {
				found = true
				continue
			}
			shortcuts = append(shortcuts, existing)
		}
		if len(shortcuts) == 0 {
			continue
		}
		entry.Shortcuts = shortcuts
		updated = append(updated, entry)
	}

	if !found {
		return fmt.Errorf("alias shortcut %q not found", shortcut)
	}

	cfg.Alias.Entries = updated
	return s.config.Save(cfg)
}

func (s *AliasService) Reset() error {
	cfg, err := s.config.Load()
	if err != nil {
		return err
	}
	cfg.Alias = domainconfig.Default().Alias
	return s.config.Save(cfg)
}

func (s *AliasService) resolveEmoji(name, explicitEmoji string) (string, error) {
	if explicitEmoji != "" {
		return explicitEmoji, nil
	}

	if match, ok := s.catalog.Lookup(name); ok {
		return match.Emoji, nil
	}

	if s.remote != nil {
		if match, ok, err := s.remote.Lookup(name); err == nil && ok {
			return match.Emoji, nil
		}
	}

	return "", fmt.Errorf("emoji not found for alias name %q", name)
}
