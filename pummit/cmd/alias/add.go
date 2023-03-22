package alias_cmd

import (
	"os"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/pummit/config"
	"github.com/HidemaruOwO/pummit/pummit/lib"
)

var isDebug bool = config.IsDebug()

func AliasAddCmd(args []string) {
	aliasAdd(args)
}

func aliasAdd(args []string) {
	aliases := lib.GetAlias()

	if len(args) == 2 {
		log.Warnf("Missing alias\n")
		os.Exit(0)
	}
	if len(args) == 3 {
		log.Warnf("Missing emoji prefix\n")
		os.Exit(0)
	}

	alias := args[2]
	prefix := args[3]
	var emoji string

	gitimojis := lib.GetGitmoji()

	for index, value := range gitimojis.Gitmojis {
		if value.Name == prefix {
			log.Debugf(isDebug, "found prefix from gitimoji\n")
			emoji = gitimojis.Gitmojis[index].Emoji
		}
	}

	if emoji == "" {
		log.Warnf("The prefix '%s' was not found in gitiemoji.\n", prefix)
		os.Exit(0)
	}

	aliases.Alias = append(aliases.Alias, []string{alias, prefix, emoji})

	lib.WriteConfig(aliases)
	log.Infof("The following has been added to the custom alias\n  %s : %s : %s\n", alias, prefix, emoji)
}
