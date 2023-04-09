package lib

import (
	"sync"

	"github.com/HidemaruOwO/nuts/log"
)

func IncludeGitimoji(prefix string) (bool, string) {
	var emoji string
	var gitimoji Gitmoji = GetGitmoji()

	var wg sync.WaitGroup

	for index, value := range gitimoji.Gitmojis {
		wg.Add(1)
		go func(index int, value Gitmojis) {
			if value.Name == prefix {
				log.Debugf(isDebug, "found prefix from gitmoji\n")
				emoji = gitimoji.Gitmojis[index].Emoji
			}
			wg.Done()
		}(index, value)
	}
	wg.Wait()

	if emoji == "" {
		return false, ""
	}
	return true, emoji
}
