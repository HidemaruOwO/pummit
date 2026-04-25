package emoji

import (
	domainemoji "github.com/HidemaruOwO/pummit/internal/domain/emoji"
	legacygitmoji "github.com/HidemaruOwO/pummit/legacy/gitmoji"
)

type GitmojiClient struct{}

func NewGitmojiClient() *GitmojiClient {
	return &GitmojiClient{}
}

func (c *GitmojiClient) Lookup(name string) (domainemoji.Match, bool, error) {
	gitmojis, err := legacygitmoji.GetAllGitmojisWithConfig(false)
	if err != nil {
		return domainemoji.Match{}, false, err
	}

	if item, ok := legacygitmoji.FindByName(name, gitmojis); ok {
		return domainemoji.Match{Name: item.Name, Emoji: item.Emoji}, true, nil
	}

	if item, ok := legacygitmoji.FindByCode(name, gitmojis); ok {
		return domainemoji.Match{Name: item.Name, Emoji: item.Emoji}, true, nil
	}

	return domainemoji.Match{}, false, nil
}
