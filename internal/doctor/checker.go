package doctor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/HidemaruOwO/pummit/internal/config"
	"github.com/HidemaruOwO/pummit/pkg/gitmoji"
)

// DiagnosticResult 診断結果を表す構造体
type DiagnosticResult struct {
	Name        string   `json:"name"`        // 診断項目名
	Status      string   `json:"status"`      // OK, WARNING, ERROR
	Message     string   `json:"message"`     // 診断結果メッセージ
	Suggestions []string `json:"suggestions"` // 問題解決の提案
}

// DiagnosticResults 診断結果のリスト
type DiagnosticResults []DiagnosticResult

// Checker 診断チェッカーのインターフェース
type Checker interface {
	Name() string
	Check() DiagnosticResult
}

// GitConfigChecker Git設定を検証するチェッカー
type GitConfigChecker struct{}

func (c *GitConfigChecker) Name() string {
	return "Git Configuration"
}

func (c *GitConfigChecker) Check() DiagnosticResult {
	result := DiagnosticResult{
		Name:        c.Name(),
		Status:      "OK",
		Message:     "Git configuration is properly set",
		Suggestions: []string{},
	}

	// Git user.nameをチェック
	nameCmd := exec.Command("git", "config", "--global", "user.name")
	nameOutput, nameErr := nameCmd.Output()

	// Git user.emailをチェック
	emailCmd := exec.Command("git", "config", "--global", "user.email")
	emailOutput, emailErr := emailCmd.Output()

	var issues []string
	if nameErr != nil || strings.TrimSpace(string(nameOutput)) == "" {
		issues = append(issues, "user.name is not configured")
		result.Suggestions = append(result.Suggestions, "git config --global user.name \"Your Name\"")
	}

	if emailErr != nil || strings.TrimSpace(string(emailOutput)) == "" {
		issues = append(issues, "user.email is not configured")
		result.Suggestions = append(result.Suggestions, "git config --global user.email \"your.email@example.com\"")
	}

	if len(issues) > 0 {
		result.Status = "ERROR"
		result.Message = fmt.Sprintf("Git configuration issues found: %s", strings.Join(issues, ", "))
	}

	return result
}

// GitRepositoryChecker Gitリポジトリの状態を確認するチェッカー
type GitRepositoryChecker struct{}

func (c *GitRepositoryChecker) Name() string {
	return "Git Repository Status"
}

func (c *GitRepositoryChecker) Check() DiagnosticResult {
	result := DiagnosticResult{
		Name:        c.Name(),
		Status:      "OK",
		Message:     "Git repository is in good state",
		Suggestions: []string{},
	}

	// .gitディレクトリの存在確認
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		result.Status = "WARNING"
		result.Message = "Current directory is not a Git repository"
		result.Suggestions = append(result.Suggestions, "Initialize repository with: git init")
		return result
	}

	// Git status確認
	statusCmd := exec.Command("git", "status", "--porcelain")
	statusOutput, err := statusCmd.Output()
	if err != nil {
		result.Status = "ERROR"
		result.Message = "Failed to check Git status"
		return result
	}

	// ステージングされたファイルの確認
	stagedFiles := 0
	unstagedFiles := 0
	lines := strings.Split(strings.TrimSpace(string(statusOutput)), "\n")
	for _, line := range lines {
		if len(line) >= 2 {
			if line[0] != ' ' && line[0] != '?' {
				stagedFiles++
			}
			if line[1] != ' ' && line[1] != '?' {
				unstagedFiles++
			}
		}
	}

	if stagedFiles == 0 && unstagedFiles > 0 {
		result.Status = "WARNING"
		result.Message = fmt.Sprintf("Found %d unstaged files", unstagedFiles)
		result.Suggestions = append(result.Suggestions, "Stage files with: git add .")
	}

	return result
}

// ConfigFileChecker 設定ファイルの妥当性を検証するチェッカー
type ConfigFileChecker struct{}

func (c *ConfigFileChecker) Name() string {
	return "Configuration Files"
}

func (c *ConfigFileChecker) Check() DiagnosticResult {
	result := DiagnosticResult{
		Name:        c.Name(),
		Status:      "OK",
		Message:     "Configuration files are valid",
		Suggestions: []string{},
	}

	// 設定ディレクトリの確認
	configDir, err := config.GetConfigDir()
	if err != nil {
		result.Status = "ERROR"
		result.Message = fmt.Sprintf("Failed to get configuration directory: %v", err)
		return result
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		result.Status = "WARNING"
		result.Message = "Configuration directory does not exist"
		result.Suggestions = append(result.Suggestions, "Run pummit once to initialize configuration")
		return result
	}

	// TOML設定ファイルの確認を優先
	tomlPath := filepath.Join(configDir, "config.toml")
	jsonPath := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(tomlPath); err == nil {
		// TOML設定ファイルの妥当性チェック
		if err := config.LoadTOMLConfig(); err != nil {
			result.Status = "ERROR"
			result.Message = fmt.Sprintf("TOML configuration file is corrupted: %v", err)
			result.Suggestions = append(result.Suggestions, "Check configuration file syntax")
			result.Suggestions = append(result.Suggestions, "Restore configuration from backup")
		}
	} else if _, err := os.Stat(jsonPath); err == nil {
		// TOMLが存在しない場合のみJSON設定ファイルをチェック
		if err := c.validateJSONConfig(jsonPath); err != nil {
			result.Status = "ERROR"
			result.Message = fmt.Sprintf("JSON configuration file is corrupted: %v", err)
			result.Suggestions = append(result.Suggestions, "Migrate to TOML format with: pummit migrate")
			result.Suggestions = append(result.Suggestions, "Restore configuration from backup")
		}
	} else {
		// 設定ファイルが全く存在しない場合
		result.Status = "WARNING"
		result.Message = "No configuration files found"
		result.Suggestions = append(result.Suggestions, "Run pummit once to initialize configuration")
	}

	return result
}

// validateJSONConfig JSON設定ファイルの妥当性を検証
func (c *ConfigFileChecker) validateJSONConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config config.Config
	return json.Unmarshal(data, &config)
}

// NetworkChecker ネットワーク接続を確認するチェッカー
type NetworkChecker struct{}

func (c *NetworkChecker) Name() string {
	return "Network Connectivity"
}

func (c *NetworkChecker) Check() DiagnosticResult {
	result := DiagnosticResult{
		Name:        c.Name(),
		Status:      "OK",
		Message:     "Network connectivity is available",
		Suggestions: []string{},
	}

	// Gitmoji APIへの接続確認
	if !gitmoji.IsOnline() {
		result.Status = "WARNING"
		result.Message = "Cannot connect to Gitmoji API (operating in offline mode)"
		result.Suggestions = append(result.Suggestions, "Check your internet connection")
		result.Suggestions = append(result.Suggestions, "Basic functionality is available in offline mode")
	}

	return result
}

// FilePermissionChecker ファイル書き込み権限を確認するチェッカー
type FilePermissionChecker struct{}

func (c *FilePermissionChecker) Name() string {
	return "File Permissions"
}

func (c *FilePermissionChecker) Check() DiagnosticResult {
	result := DiagnosticResult{
		Name:        c.Name(),
		Status:      "OK",
		Message:     "File permissions are properly configured",
		Suggestions: []string{},
	}

	// 設定ディレクトリの書き込み権限確認
	configDir, err := config.GetConfigDir()
	if err != nil {
		result.Status = "ERROR"
		result.Message = fmt.Sprintf("Failed to get configuration directory: %v", err)
		return result
	}

	// ディレクトリが存在しない場合は作成を試行
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			result.Status = "ERROR"
			result.Message = fmt.Sprintf("Failed to create configuration directory: %v", err)
			result.Suggestions = append(result.Suggestions, "Check directory permissions")
			return result
		}
	}

	// テストファイルの作成と削除で書き込み権限確認
	testFile := filepath.Join(configDir, ".permission_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		result.Status = "ERROR"
		result.Message = fmt.Sprintf("No write permission to configuration directory: %v", err)
		result.Suggestions = append(result.Suggestions, "Check directory permissions")
		return result
	}

	// テストファイルを削除
	if err := os.Remove(testFile); err != nil {
		result.Status = "ERROR"
		result.Message = fmt.Sprintf("Failed to remove test file: %v", err)
		result.Suggestions = append(result.Suggestions, "Check directory permissions")
		return result
	}

	return result
}

// RunAllChecks すべての診断チェックを実行
func RunAllChecks() DiagnosticResults {
	checkers := []Checker{
		&GitConfigChecker{},
		&GitRepositoryChecker{},
		&ConfigFileChecker{},
		&NetworkChecker{},
		&FilePermissionChecker{},
	}

	var results DiagnosticResults
	for _, checker := range checkers {
		results = append(results, checker.Check())
	}

	return results
}

// HasErrors 診断結果にエラーが含まれているかチェック
func (results DiagnosticResults) HasErrors() bool {
	for _, result := range results {
		if result.Status == "ERROR" {
			return true
		}
	}
	return false
}

// HasWarnings 診断結果に警告が含まれているかチェック
func (results DiagnosticResults) HasWarnings() bool {
	for _, result := range results {
		if result.Status == "WARNING" {
			return true
		}
	}
	return false
}
