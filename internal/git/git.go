package git

import (
	"os"
	"os/exec"
	"strings"
)

type CommitMessage struct {
	Emoji   string
	Message string
}

// カレントディレクトリがGitリポジトリか
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// コミットを作成
func Commit(cm CommitMessage) error {
	message := cm.Emoji + " " + cm.Message
	cmd := exec.Command("git", "commit", "-m", message)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// 今のブランチ名の表示
func GetBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
