package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/alias"
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

// コミットを作成
func Commit(cm CommitMessage) error {
	log := logger.New()

	// 変更済みのファイルを取得
	changed, err := GetChangedFiles()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// 絵文字エイリアス実装
	found, emoji := alias.GetEmoji(cm.Emoji)
	if config.CurrentConfig.UseAlias {
		if found {
			// 55行目のTODOを実装したらここはemojiがprefixになる
			// cm.Emoji = fmt.Sprintf(":%s:", prefix)
			cm.Emoji = emoji
		} else {
			// エイリアスが見つからない場合はそのまま
			cm.Emoji = fmt.Sprintf(":%s:", cm.Emoji)
		}

	} else {
		// エイリアスが見つからない場合は"::"で括る
		cm.Emoji = fmt.Sprintf(":%s:", cm.Emoji)
	}

	if config.CurrentConfig.UseRawEmoji {
		// TODO
		// 絵文字には変換しない場合なのでalias.findAlias関数にemoji prefixを返すようにも実装してあげるようにする必要がある
		// (e.g) found, emoji, prefix := alias.GetEmoji(cm.Emoji)

		if found {
			cm.Emoji = emoji
		} else {
			// エイリアスが見つからない場合はそのまま
			cm.Emoji = fmt.Sprintf(":%s:", cm.Emoji)
		}
	}

	// コミットメッセージが長くなりすぎないようにするため
	if config.CurrentConfig.UseFilesLength {
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
