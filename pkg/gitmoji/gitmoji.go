package gitmoji

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type Gitmoji struct {
	Emoji       string `json:"emoji"`
	Entity      string `json:"entity"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type GitmojiResponse struct {
	Gitmojis []Gitmoji `json:"gitmojis"`
}

// maxGitmojiResponseSize prevents memory exhaustion from large payloads.
// Limit is set to 1MB (1024 * 1024 bytes), a reasonable upper bound.
const maxGitmojiResponseSize = 1024 * 1024

var (
	// ErrNonOK allows callers to detect unexpected HTTP status codes.
	ErrNonOK = errors.New("non-200 from server")

	// ErrDecode helps distinguish malformed payloads from transport failures.
	ErrDecode = errors.New("failed to decode response")
)

// Fetch retrieves gitmojis from the given URL using the provided client.
func Fetch(ctx context.Context, client *http.Client, url string) (
	GitmojiResponse, error,
) {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return GitmojiResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return GitmojiResponse{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return GitmojiResponse{}, fmt.Errorf(
			"%w: status=%d", ErrNonOK, resp.StatusCode,
		)
	}

	var result GitmojiResponse
	r := io.LimitReader(resp.Body, maxGitmojiResponseSize)
	dec := json.NewDecoder(r)
	if err := dec.Decode(&result); err != nil {
		return GitmojiResponse{}, fmt.Errorf("%w: %v", ErrDecode, err)
	}

	return result, nil
}

// オフライン対応のGitmoji取得（フォールバック機能付き）
func GetAllGitmojis() ([]Gitmoji, error) {
	return GetAllGitmojisWithConfig(false)
}

// 設定可能なGitmoji取得関数
func GetAllGitmojisWithConfig(offlineMode bool) ([]Gitmoji, error) {
	if offlineMode {
		// オフラインモードでは即座にエラーを返してフォールバックを促す
		return nil, errors.New("offline mode is enabled")
	}

	// オフライン状態検出のためのタイムアウトを短縮（3秒）
	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 1 * time.Second, // 接続タイムアウト
			}).DialContext,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://raw.githubusercontent.com/carloscuesta/gitmoji/master/packages/gitmojis/src/gitmojis.json", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		// ネットワークエラーの場合、より具体的なエラーメッセージを返す
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, errors.New("network connection timeout (possibly offline)")
		}
		return nil, errors.New("network connection failed")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch gitmojis from API")
	}

	var response GitmojiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Gitmojis, nil
}

// ネットワーク状態の簡易チェック
func IsOnline() bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 1 * time.Second,
			}).DialContext,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 軽量なGETリクエストでネットワーク状態をチェック
	req, err := http.NewRequestWithContext(ctx, "HEAD", "https://api.github.com", nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode == http.StatusOK
}

func FindByCode(code string, gitmojis []Gitmoji) (Gitmoji, bool) {
	for _, gitmoji := range gitmojis {
		if gitmoji.Code == code {
			return gitmoji, true
		}
	}
	return Gitmoji{}, false
}

func FindByName(name string, gitmojis []Gitmoji) (Gitmoji, bool) {
	for _, gitmoji := range gitmojis {
		if gitmoji.Name == name {
			return gitmoji, true
		}
	}
	return Gitmoji{}, false
}
