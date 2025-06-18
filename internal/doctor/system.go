package doctor

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/variable"
)

// SystemInfo システム情報を格納する構造体
type SystemInfo struct {
	OS            string // オペレーティングシステム
	Architecture  string // アーキテクチャ
	GoVersion     string // Go バージョン
	PummitVersion string // Pummit バージョン
	GitVersion    string // Git バージョン
}

// GetSystemInfo システム情報を収集して返す
func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:            runtime.GOOS,
		Architecture:  runtime.GOARCH,
		PummitVersion: variable.VERSION,
		GoVersion:     getGoVersion(),
		GitVersion:    getGitVersion(),
	}

	return info
}

// getGoVersion Goのバージョンを取得
func getGoVersion() string {
	version := runtime.Version()
	// "go1.21.0" -> "1.21.0" の形式に変換
	if strings.HasPrefix(version, "go") {
		return strings.TrimPrefix(version, "go")
	}
	return version
}

// getGitVersion Gitのバージョンを取得
func getGitVersion() string {
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "Not installed or not accessible"
	}

	// "git version 2.39.0" -> "2.39.0" の形式に変換
	versionStr := strings.TrimSpace(string(output))
	parts := strings.Fields(versionStr)
	if len(parts) >= 3 {
		return parts[2]
	}

	return versionStr
}

// FormatSystemInfo システム情報を人間が読みやすい形式でフォーマット
func (info SystemInfo) FormatSystemInfo() string {
	var builder strings.Builder

	builder.WriteString("System Information:\n")
	builder.WriteString(fmt.Sprintf("  OS: %s %s\n", info.OS, info.Architecture))
	builder.WriteString(fmt.Sprintf("  Go Version: %s\n", info.GoVersion))
	builder.WriteString(fmt.Sprintf("  Pummit Version: %s\n", info.PummitVersion))
	builder.WriteString(fmt.Sprintf("  Git Version: %s\n", info.GitVersion))

	return builder.String()
}
