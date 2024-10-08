package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/src/config"
	"github.com/HidemaruOwO/pummit/src/lib"
)

func RootCmd(prefix string, subject string) {
	gitCommit(prefix, subject)
}

func gitCommit(prefix string, subject string) {
	gitChangeCmd := exec.Command("git", "diff", "--name-only", "--cached")
	gitChangeOutput, err := gitChangeCmd.Output()
	log.Debugf(config.IS_DEBUG, "%s", gitChangeOutput)
	if err != nil {
		log.ErrorExit(err)
	}
	gitChange := strings.TrimSpace(string(gitChangeOutput))

	if gitChange == "" {
		log.Infof("checked the Git diff, but it looks like there are no new changes\n")
		os.Exit(0)
	}

	args := os.Args[1:]

	if len(os.Args) == 2 {
		arr := strings.Split(args[0], " ")
		if len(arr) == 1 {
			log.Infof("The emoji prefix is missing. Please add it as the first argument\n")
			return
		}
		prefix = arr[0]
		subject = strings.Join(arr[1:], " ")
	} else {
		prefix = args[0]
		subject = args[1]
	}

	aliasList := lib.GetAliasList()

	var wg sync.WaitGroup
	var mu sync.Mutex

	for index, value := range aliasList {
		wg.Add(1)
		go func(index int, value []string) {
			defer wg.Done()
			if value[0] == prefix {
				log.Debugf(config.IS_DEBUG, "Found prefix %s\n", aliasList[index][1])
				mu.Lock()
				prefix = aliasList[index][1]
				mu.Unlock()
			}
		}(index, value)
	}
	wg.Wait()

	_, emoji := lib.IncludeGitimoji(prefix)
	appConfig := lib.GetAppConfig()

	gitChange = strings.ReplaceAll(gitChange, "\n", ", ")

	var commigMessage string

	var message string

	if *appConfig.WriteEmoji {
		if emoji != "" {
			prefix = emoji
		} else {
			prefix = fmt.Sprintf(":%s:", prefix)
		}
	}

	if *appConfig.UseLimitPathesLength {
		if *appConfig.LimitPathesLength < len([]rune(gitChange)) {
			gitChange = fmt.Sprintf("%s...", gitChange[:*appConfig.LimitPathesLength])
		}

		message = fmt.Sprintf("%s (%s)", subject, gitChange)
	} else {
		message = fmt.Sprintf("%s", subject)
	}

	commigMessage = fmt.Sprintf("%s %s", prefix, message)

	cmd := exec.Command("git", "commit", "-m", commigMessage)
	_, err = cmd.Output()
	if err != nil {
		log.ErrorExit(err)
	}
}
