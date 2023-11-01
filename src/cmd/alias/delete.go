package alias_cmd

import (
	"sync"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/src/config"
	"github.com/HidemaruOwO/pummit/src/lib"
)

func AliasDeleteCmd(targets []string) {
	aliasDelete(targets)
}

func aliasDelete(targets []string) {
	alias := lib.GetAlias()
	slice := alias.Alias

	log.Debugf(config.IS_DEBUG, "Removing %v\n", targets)

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
