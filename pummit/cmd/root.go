package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/src/lib"
)

func RootCmd() {
	gitCommit()
}

func gitCommit() {
	var emoji string
	var subject string

	gitChangeCmd := exec.Command("git", "diff", "--name-only", "--cached", "HEAD")
	gitChangeOutput, err := gitChangeCmd.Output()
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
		emoji = arr[0]
		subject = strings.Join(arr[1:], "")
	} else {
		emoji = args[0]
		subject = args[1]
	}

	array := lib.GetAliasList()
	ch := make(chan int)
	for i, value := range array {
		go func(index int, value []string) {
			if value[0][0] == emoji[0] {
				ch <- index
			} else {
				ch <- -1
			}
		}(i, value)
	}

	// 結果の収集
	var num int
	for i := 0; i < len(array); i++ {
		res := <-ch
		if res >= 0 {
			num = res
			break
		}
	}

	emoji = array[num][1]

	gitChange = strings.ReplaceAll(gitChange, "\n", ", ")

	commitMsg := fmt.Sprintf(":%s: %s (%s)", emoji, subject, gitChange)
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	_, err = cmd.Output()
	if err != nil {
		log.ErrorExit(err)
	}
}
