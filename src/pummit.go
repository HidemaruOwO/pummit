package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func gitCommit() {
	var emoji string
	var subject string
	gitChangeCmd := exec.Command("git", "diff", "--name-only", "--cached", "HEAD")
	gitChangeOutput, err := gitChangeCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	gitChange := strings.TrimSpace(string(gitChangeOutput))

	if gitChange == "" {
		fmt.Println("No change")
		return
	}
	if len(os.Args) == 1 {
		fmt.Println("No args")
		return
	}

	args := os.Args[1:]

	if len(os.Args) == 2 {
		arr := strings.Split(args[0], " ")
		if len(arr) == 1 {
			fmt.Println("No emoji prefix")
			return
		}
		emoji = arr[0]
		subject = strings.Join(arr[1:], "")
	} else {
		emoji = args[0]
		subject = args[1]
	}

	gitChange = strings.ReplaceAll(gitChange, "\n", ", ")

	commitMsg := fmt.Sprintf(":%s: %s (%s)", emoji, subject, gitChange)
	cmd := exec.Command("git", "commit", "-m", commitMsg)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	gitCommit()
}
