package emojis

import (
	"encoding/json"
	"fmt"

	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/HidemaruOwO/pummit/pkg/gitmoji"
)

type Emoji struct {
	Emoji string `json:"emoji"`
	Name  string `json:"name"`
}

var (
	mapping   map[string]string
	allEmojis []Emoji
)

func init() {
	mapping = make(map[string]string)

	err := json.Unmarshal([]byte(variable.EMOJIS_JSON), &allEmojis)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse emojis.json: %v", err))
	}

	for _, e := range allEmojis {
		mapping[e.Name] = e.Emoji
	}
}

// オフライン対応の絵文字取得（フォールバック機能付き）
func GetEmojiByName(name string) (string, error) {
	return GetEmojiByNameWithFallback(name, false)
}

// オフラインフラグ指定での絵文字取得
func GetEmojiByNameOffline(name string, offlineMode bool) (string, error) {
	return GetEmojiByNameWithFallback(name, offlineMode)
}

// フォールバック機能付き絵文字取得
func GetEmojiByNameWithFallback(name string, offlineMode bool) (string, error) {
	// まず埋め込みデータから検索
	emoji, ok := mapping[name]
	if ok {
		return emoji, nil
	}

	// オフラインモードの場合は埋め込みデータのみを使用
	if offlineMode {
		return "", fmt.Errorf("emoji not found for name: %s (offline mode)", name)
	}

	// オンラインの場合はGitmoji APIから取得を試行
	gitmojis, err := gitmoji.GetAllGitmojisWithConfig(false)
	if err != nil {
		// API呼び出しに失敗した場合、埋め込みデータで再試行
		return "", fmt.Errorf("emoji not found for name: %s (API unavailable: %v)", name, err)
	}

	// Gitmoji APIから取得したデータで検索
	gitmojiEmoji, found := gitmoji.FindByName(name, gitmojis)
	if found {
		return gitmojiEmoji.Emoji, nil
	}

	// コード名でも検索
	gitmojiEmoji, found = gitmoji.FindByCode(name, gitmojis)
	if found {
		return gitmojiEmoji.Emoji, nil
	}

	return "", fmt.Errorf("emoji not found for name: %s", name)
}

func GetAllEmojis() map[string]string {
	copiedMap := make(map[string]string, len(mapping))
	for k, v := range mapping {
		copiedMap[k] = v
	}
	return copiedMap
}

func GetAllEmojiStructs() []Emoji {
	return allEmojis
}
