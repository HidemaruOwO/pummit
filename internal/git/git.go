package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/internal/emojis"
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
	return CommitWithOfflineMode(cm, false)
}

// オフラインモード対応のコミット関数
func CommitWithOfflineMode(cm CommitMessage, offlineMode bool) error {
	return commitWithOfflineModeInternal(cm, offlineMode)
}

// オーバーロード: 文字列パラメータでのコミット関数（MCP用）
func CommitWithOfflineModeStrings(emoji, message string, offlineMode bool) error {
	cm := CommitMessage{
		Emoji:   emoji,
		Message: message,
	}
	return commitWithOfflineModeInternal(cm, offlineMode)
}

// 内部実装
func commitWithOfflineModeInternal(cm CommitMessage, offlineMode bool) error {
	log := logger.New()

	// 変更済みのファイルを取得
	changed, err := GetChangedFiles()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	enteredEmoji := cm.Emoji
	found, prefix, emoji := alias.GetEmoji(enteredEmoji)

	// TOMLベースの設定を使用
	useRawEmoji := config.CurrentTOMLConfig.Base.Emoji
	useAlias := config.CurrentTOMLConfig.Alias.Enabled

	if useRawEmoji {
		// 絵文字モード
		if useAlias && found {
			cm.Emoji = emoji
		} else {
			cm.Emoji = ConvertToEmojiWithOfflineMode(enteredEmoji, offlineMode)
		}
	} else {
		// :name: モード
		if useAlias && found {
			cm.Emoji = fmt.Sprintf(":%s:", prefix)
		} else {
			cm.Emoji = fmt.Sprintf(":%s:", enteredEmoji)
		}
	}

	// コミットメッセージが長くなりすぎないようにするため
	filesLength := config.CurrentTOMLConfig.Base.FilesLength
	if filesLength > 0 && filesLength < len([]rune(changed)) {
		changed = fmt.Sprintf("%s...", changed[:filesLength])
	}

	// message := cm.Emoji + " " + cm.Message + " " + changed
	message := fmt.Sprintf("%s %s (%s)", cm.Emoji, cm.Message, changed)
	cmd := exec.Command("git", "commit", "-m", message)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// AddFiles 指定されたファイルをステージングする
func AddFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

// GetStagedFiles ステージされたファイル一覧を取得する
func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	staged := strings.TrimSpace(string(output))
	if staged == "" {
		return []string{}, nil
	}

	return strings.Split(staged, "\n"), nil
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

// GetChangedFilesList 変更されたファイル一覧をスライスで取得する
func GetChangedFilesList() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	changed := strings.TrimSpace(string(output))
	if changed == "" {
		return []string{}, nil
	}

	return strings.Split(changed, "\n"), nil
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

func ConvertToEmoji(name string) string {
	return ConvertToEmojiWithOfflineMode(name, false)
}

// オフラインモード対応の絵文字変換
func ConvertToEmojiWithOfflineMode(name string, offlineMode bool) string {
	emoji, err := emojis.GetEmojiByNameOffline(name, offlineMode)
	if err != nil {
		// エラーが発生した場合は名前をコロンで囲んで返す
		return fmt.Sprintf(":%s:", name)
	}
	return emoji
}
