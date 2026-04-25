package emoji

import (
	"encoding/json"
	"sync"

	domainemoji "github.com/HidemaruOwO/pummit/internal/domain/emoji"
	"github.com/HidemaruOwO/pummit/internal/variable"
)

type Catalog struct {
	once    sync.Once
	mapping map[string]string
	err     error
}

func NewCatalog() *Catalog {
	return &Catalog{}
}

func (c *Catalog) Lookup(name string) (domainemoji.Match, bool) {
	if err := c.load(); err != nil {
		return domainemoji.Match{}, false
	}

	emoji, ok := c.mapping[name]
	if !ok {
		return domainemoji.Match{}, false
	}

	return domainemoji.Match{Name: name, Emoji: emoji}, true
}

func (c *Catalog) load() error {
	c.once.Do(func() {
		var items []domainemoji.Match
		if err := json.Unmarshal([]byte(variable.EMOJIS_JSON), &items); err != nil {
			c.err = err
			return
		}

		c.mapping = make(map[string]string, len(items))
		for _, item := range items {
			c.mapping[item.Name] = item.Emoji
		}
	})

	return c.err
}
