package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/pummit/config"
	"github.com/HidemaruOwO/pummit/pummit/lib"
)

var isDebug bool = config.IsDebug()

func RootCmd(prefix string, subject string) {
	gitCommit(prefix, subject)
}

func gitCommit(prefix string, subject string) {
	gitChangeCmd := exec.Command("git", "diff", "--name-only", "--cached")
	gitChangeOutput, err := gitChangeCmd.Output()
	log.Debugf(isDebug, "%s", gitChangeOutput)
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

	array := lib.GetAliasList()

	var wg sync.WaitGroup
	var mu sync.Mutex
	for index, value := range array {
		wg.Add(1)
		go func(index int, value []string) {
			defer wg.Done()
			if value[0] == prefix {
				log.Debugf(isDebug, "Found prefix %s\n", array[index][1])
				mu.Lock()
				prefix = array[index][1]
				mu.Unlock()
			}
		}(index, value)
	}
	wg.Wait()

	gitChange = strings.ReplaceAll(gitChange, "\n", ", ")

	commitMsg := fmt.Sprintf(":%s: %s (%s)", prefix, subject, gitChange)
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	_, err = cmd.Output()
	if err != nil {
		log.ErrorExit(err)
	}
}
