package alias_cmd

import (
	"os"
	"sync"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/pummit/lib"
)

func AliasDeleteCmd(args []string) {
	aliasDelete(args)
}

func aliasDelete(args []string) {
	alias := lib.GetAlias()

	if len(args) == 2 {
		log.Warnf("The alias name to be deleted must be added as the third argument\n")
		os.Exit(0)
	} else {
		slice := alias.Alias
		targets := args[2:]
		log.Debugf(isDebug, "Removing %v\n", targets)

		var wg sync.WaitGroup
		wg.Add(len(targets))
		for _, t := range targets {
			go func(target string) {
				defer wg.Done()
				slice = removeSlice(slice, target)
				log.Infof("Removed %s alias\n", target)
			}(t)
		}
		wg.Wait()

		alias.Alias = slice
		lib.WriteConfig(alias)
	}
}

func removeSlice(slice [][]string, target string) [][]string {
	var result [][]string
	for _, s := range slice {
		found := false
		for _, v := range s {
			if v == target {
				found = true
				break
			}
		}
		if !found {
			result = append(result, s)
		}
	}
	return result
}
