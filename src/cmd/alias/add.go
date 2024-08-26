package alias_cmd

import (
	"os"
	"strings"
	"sync"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/src/lib"
)

func AliasAddCmd(alias string, prefix string) {
	aliasAdd(alias, prefix)
}

func aliasAdd(alias string, prefix string) {
	aliases := lib.GetAppConfig()

	if alias != strings.TrimSpace(alias) || prefix != strings.TrimSpace(prefix) {
		log.Warnf("Please do not include spaces in alias and emoji prefixes\n")
		os.Exit(0)
	}

	isGitimoji, emoji := lib.IncludeGitimoji(prefix)

	aliasList := lib.GetAliasList()

	var wg sync.WaitGroup
	for index, value := range aliasList {
		wg.Add(1)
		go func(index int, value []string) {
			if value[0] == alias {
				log.Warnf("An alias with the same name already exists\n")
				os.Exit(0)
			}
			wg.Done()
		}(index, value)
	}
	wg.Wait()

	if !isGitimoji {
		log.Warnf("The prefix '%s' was not found in gitiemoji.\n", prefix)
		os.Exit(0)
	}

	*aliases.Alias = append(*aliases.Alias, []string{alias, prefix, emoji})

	lib.WriteConfig(aliases)
	log.Infof(
		"The following has been added to the custom alias\n  %s : %s : %s\n",
		alias,
		prefix,
		emoji,
	)
}
