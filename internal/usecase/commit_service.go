package usecase

import (
	"fmt"
	"regexp"
	"strings"

	domaincommit "github.com/HidemaruOwO/pummit/internal/domain/commit"
	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
	domainemoji "github.com/HidemaruOwO/pummit/internal/domain/emoji"
)

type ConfigLoader interface {
	Load() (domainconfig.Config, error)
}

type GitClient interface {
	StagedFiles() ([]string, error)
	CurrentBranch() (string, error)
	Commit(message string) error
}

type Catalog interface {
	Lookup(name string) (domainemoji.Match, bool)
}

type GitmojiLookup interface {
	Lookup(name string) (domainemoji.Match, bool, error)
}

type CommitService struct {
	config  ConfigLoader
	git     GitClient
	catalog Catalog
	remote  GitmojiLookup
}

func NewCommitService(config ConfigLoader, git GitClient, catalog Catalog, remote GitmojiLookup) *CommitService {
	return &CommitService{config: config, git: git, catalog: catalog, remote: remote}
}

func (s *CommitService) CommitExplicit(emojiInput, message string) (domaincommit.Result, error) {
	return s.commit(domaincommit.Request{EmojiInput: emojiInput, Message: message})
}

func (s *CommitService) CommitAuto(message string) (domaincommit.Result, error) {
	return s.commit(domaincommit.Request{Message: message, AutoEmoji: true})
}

func (s *CommitService) commit(req domaincommit.Request) (domaincommit.Result, error) {
	cfg, err := s.config.Load()
	if err != nil {
		return domaincommit.Result{}, err
	}

	stagedFiles, err := s.git.StagedFiles()
	if err != nil {
		return domaincommit.Result{}, err
	}

	if len(stagedFiles) == 0 {
		return domaincommit.Result{}, nil
	}

	prefix, err := s.resolvePrefix(req, cfg)
	if err != nil {
		return domaincommit.Result{}, err
	}

	files := strings.Join(stagedFiles, ", ")
	if cfg.Base.FilesLength > 0 {
		files = domaincommit.TruncateFiles(files, cfg.Base.FilesLength)
	}

	formatted := domaincommit.FormatMessage(prefix, req.Message, files)
	if err := s.git.Commit(formatted); err != nil {
		return domaincommit.Result{}, err
	}

	return domaincommit.Result{Committed: true, Message: formatted}, nil
}

func (s *CommitService) resolvePrefix(req domaincommit.Request, cfg domainconfig.Config) (string, error) {
	if req.AutoEmoji {
		name, err := s.resolveAutoName(cfg)
		if err != nil {
			return "", err
		}
		return s.resolveExplicitPrefix(name, cfg), nil
	}

	return s.resolveExplicitPrefix(req.EmojiInput, cfg), nil
}

func (s *CommitService) resolveAutoName(cfg domainconfig.Config) (string, error) {
	branch, err := s.git.CurrentBranch()
	if err != nil {
		return "", err
	}

	if cfg.BranchMapping.Enabled {
		for _, rule := range cfg.BranchMapping.Rules {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				return "", err
			}
			if re.MatchString(branch) {
				return rule.Emoji, nil
			}
		}

		if cfg.BranchMapping.Fallback != "" {
			return cfg.BranchMapping.Fallback, nil
		}
	}

	return "construction", nil
}

func (s *CommitService) resolveExplicitPrefix(input string, cfg domainconfig.Config) string {
	if cfg.Alias.Enabled {
		if match, ok := lookupAlias(cfg.Alias.Entries, input); ok {
			return formatPrefix(match.Name, match.Emoji, cfg.Base.Emoji)
		}
	}

	if match, ok := s.catalog.Lookup(input); ok {
		return formatPrefix(match.Name, match.Emoji, cfg.Base.Emoji)
	}

	if s.remote != nil {
		if match, ok, err := s.remote.Lookup(input); err == nil && ok {
			return formatPrefix(match.Name, match.Emoji, cfg.Base.Emoji)
		}
	}

	return fmt.Sprintf(":%s:", input)
}

func formatPrefix(name, emoji string, useRaw bool) string {
	if useRaw {
		return emoji
	}

	return fmt.Sprintf(":%s:", name)
}

func lookupAlias(entries []domainconfig.AliasEntry, input string) (domainemoji.Match, bool) {
	for _, entry := range entries {
		if entry.Name == input {
			return domainemoji.Match{Name: entry.Name, Emoji: entry.Emoji}, true
		}

		for _, shortcut := range entry.Shortcuts {
			if shortcut == input {
				return domainemoji.Match{Name: entry.Name, Emoji: entry.Emoji}, true
			}
		}
	}

	return domainemoji.Match{}, false
}
