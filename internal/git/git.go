package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/alias"
	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/internal/emojis"
	"github.com/HidemaruOwO/pummit/internal/logger"
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
	if filesLength > 0 {
		changed = truncateWithEllipsis(changed, filesLength)
	}

	// message := cm.Emoji + " " + cm.Message + " " + changed
	message := fmt.Sprintf("%s %s (%s)", cm.Emoji, cm.Message, changed)
	cmd := exec.Command("git", "commit", "-m", message)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// マルチバイト文字を含むファイル名を安全に切り詰める
func truncateWithEllipsis(s string, limit int) string {
	if limit <= 0 {
		return s
	}

	runes := []rune(s)
	if limit < len(runes) {
		return fmt.Sprintf("%s...", string(runes[:limit]))
	}

	return s
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

// GetChangedFilesList は `git status --porcelain=v1 -u` の結果を解析して
// 変更(新規/更新/削除/リネーム) されたファイルパスを返す。
// 例外時は error を返す。
func GetChangedFilesList() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain=v1", "-u")
	rawOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 改行コードだけを除去（前後の空白は残す）
	rawOut = bytes.TrimRight(rawOut, "\n\r")

	scanner := bufio.NewScanner(bytes.NewReader(rawOut))
	files := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 { // "XY␠" で最低 3byte + filepath
			continue
		}

		// 3byte 目(インデックス 3) からがパス
		path := line[3:]

		// リネームの場合: "R  old -> new"
		if strings.Contains(path, " -> ") {
			parts := strings.SplitN(path, " -> ", 2)
			if len(parts) == 2 {
				path = parts[1]
			}
		}

		path = strings.TrimSpace(path)
		if path != "" {
			files = append(files, path)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return files, nil
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
