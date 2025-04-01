package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/pkg/logger"
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

// TODO これをroot.goで実行するように実装する。
// コミットを作成
func Commit(cm CommitMessage) error {
	log := logger.New()

	// TODO alias.goを直してから実装する、
	// 参考: app/cmd/root.go:76
	// cm.EmojiをUseRawEmojiに従ってプレフィックスから絵文字に置き換える
	if config.CurrentConfig.UseRawEmoji {

	}

	// 変更済みのファイルを取得
	changed, err := GetChangedFiles()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	if config.CurrentConfig.UseFilesLength {
		// コミットメッセージが長くなりすぎないようにするため
		if config.CurrentConfig.FilesLength < len([]rune(changed)) {
			changed = fmt.Sprintf("%s...", changed[:config.CurrentConfig.FilesLength])
		}
	}

	// message := cm.Emoji + " " + cm.Message + " " + changed
	message := fmt.Sprintf("%s %s (%s)", cm.Emoji, cm.Message, changed)
	cmd := exec.Command("git", "commit", "-m", message)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func GetChangedFiles() (string, error) {
	log := logger.New()

	// 変更されたファイルを確認
	cmd := exec.Command("git", "diff", "--name-only", "--cached")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	changed := strings.TrimSpace(string(stdout))

	if changed == "" {
		log.Info("Nothing to changed files.")
		os.Exit(0)
	}

	return strings.ReplaceAll(changed, "\n", ", "), nil
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
