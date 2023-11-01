package lib

import (
	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/src/config"
	"sync"
)

func IncludeGitimoji(prefix string) (bool, string) {
	var emoji string
	var gitimoji Gitmoji = GetGitmoji()

	log.Debugf(config.IS_DEBUG, "prefix: %s\n", prefix)

	var wg sync.WaitGroup

	for index, value := range gitimoji.Gitmojis {
		wg.Add(1)
		go func(index int, value Gitmojis) {
			if value.Name == prefix {
				log.Debugf(config.IS_DEBUG, "found prefix from gitmoji\n")
				emoji = gitimoji.Gitmojis[index].Emoji
			}
			wg.Done()
		}(index, value)
	}
	wg.Wait()

	log.Debugf(config.IS_DEBUG, "gitimoji: %s\n", emoji)

	if emoji == "" {
		return false, ""
	}
	return true, emoji
}
