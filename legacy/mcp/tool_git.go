package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/HidemaruOwO/pummit/legacy/git"
)

// registerGitTools Git関連のツールをMCPサーバーに登録する
func registerGitTools(s *server.MCPServer) error {
	// git.smart_commit ツール: LLM自律ワークフローの中核ツール
	smartCommitTool := mcp.NewTool("git.smart_commit",
		mcp.WithDescription("PRIMARY TOOL FOR COMMITTING. Intelligently analyze repository changes and create a commit with auto-generated message and emoji. Use this when the user simply says 'commit', 'make a commit', or gives minimal instructions like 'commit the changes' without specifying exact files or messages. Handles the complete workflow: file analysis, message generation, and execution."),
		mcp.WithString("message_hint", mcp.Description("Optional hint for commit message")),
		mcp.WithBoolean("auto_approve", mcp.Description("Auto approve commit without confirmation (default: false)")),
	)
	s.AddTool(smartCommitTool, handleGitSmartCommit)

	// git.status ツール: Gitリポジトリの状態確認
	statusTool := mcp.NewTool("git.status",
		mcp.WithDescription("Get the current status of the Git repository, including staged and unstaged files."),
	)
	s.AddTool(statusTool, handleGitStatus)

	// git.commit ツール: 詳細指定用コミット
	commitTool := mcp.NewTool("git.commit",
		mcp.WithDescription("Create a Git commit with specific user-provided emoji, message, and staging preferences. Use this ONLY when the user provides explicit commit details (specific message, emoji, or file selection). For simple 'commit' requests, use git.smart_commit instead."),
		mcp.WithString("emoji", mcp.Required(), mcp.Description("Emoji name or alias to use for the commit")),
		mcp.WithString("message", mcp.Required(), mcp.Description("Commit message")),
		mcp.WithBoolean("offline", mcp.Description("Use offline mode (default: false)")),
	)
	s.AddTool(commitTool, handleGitCommit)

	// git.add_files ツール: ファイルのステージング
	addTool := mcp.NewTool("git.add_files",
		mcp.WithDescription("Stage files for commit. If no files are specified, stages all modified files."),
		mcp.WithArray("files", mcp.Description("Array of file paths to stage (optional)")),
	)
	s.AddTool(addTool, handleGitAddFiles)

	// git.get_edited_files ツール: 変更ファイル一覧取得
	getEditedFilesTool := mcp.NewTool("git.get_edited_files",
		mcp.WithDescription("Get a list of files that have been modified, added, or deleted. Returns files in short format like 'git status -s'."),
	)
	s.AddTool(getEditedFilesTool, handleGetEditedFiles)

	// git.get_current_branch ツール: 現在のブランチ取得
	getCurrentBranchTool := mcp.NewTool("git.get_current_branch",
		mcp.WithDescription("Get the name of the current Git branch."),
	)
	s.AddTool(getCurrentBranchTool, handleGetCurrentBranch)

	return nil
}

// handleGitSmartCommit スマートコミット機能を処理する（メイン機能）
func handleGitSmartCommit(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 現状把握フェーズ
	branch, err := git.GetBranch()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current branch: %v", err)), nil
	}

	changedFiles, err := git.GetChangedFilesList()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get changed files: %v", err)), nil
	}

	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get staged files: %v", err)), nil
	}

	// 引数の取得
	args := request.GetArguments()

	// メッセージヒントの取得
	messageHint := ""
	if hint, exists := args["message_hint"]; exists {
		if hintStr, ok := hint.(string); ok {
			messageHint = hintStr
		}
	}

	// 自動承認フラグの取得
	autoApprove := false
	if approve, exists := args["auto_approve"]; exists {
		if approveBool, ok := approve.(bool); ok {
			autoApprove = approveBool
		}
	}

	// ファイルがステージされていない場合、全ての変更ファイルをステージング
	if len(stagedFiles) == 0 && len(changedFiles) > 0 {
		err = git.AddFiles(changedFiles)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to stage files: %v", err)), nil
		}
		stagedFiles = changedFiles
	}

	// コミットする内容がない場合
	if len(stagedFiles) == 0 {
		return mcp.NewToolResultText("No changes to commit. Repository is clean."), nil
	}

	// ブランチ名から絵文字を推測
	suggestedEmoji := suggestEmojiFromBranch(branch)

	// スマートなコミットメッセージ生成
	commitMessage := generateSmartCommitMessage(messageHint, stagedFiles, branch)

	// 承認が必要な場合の計画表示
	if !autoApprove {
		return mcp.NewToolResultText(fmt.Sprintf("Commit plan ready:\nBranch: %s\nStaged files: %v\nSuggested emoji: %s\nCommit message: %s\nFull commit: %s %s\n\nTo execute, call git.commit with emoji='%s' and message='%s'",
			branch, stagedFiles, suggestedEmoji, commitMessage, suggestedEmoji, commitMessage, suggestedEmoji, commitMessage)), nil
	}

	// 自動承認の場合は直接コミット実行
	cm := git.CommitMessage{
		Emoji:   suggestedEmoji,
		Message: commitMessage,
	}

	err = git.CommitWithOfflineMode(cm, false)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to commit: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully committed: %s %s", suggestedEmoji, commitMessage)), nil
}

// suggestEmojiFromBranch ブランチ名から絵文字を推測する
func suggestEmojiFromBranch(branch string) string {
	branch = strings.ToLower(branch)

	if strings.Contains(branch, "feature/") || strings.Contains(branch, "feat/") {
		return "sparkles"
	}
	if strings.Contains(branch, "fix/") || strings.Contains(branch, "bugfix/") || strings.Contains(branch, "hotfix/") {
		return "bug"
	}
	if strings.Contains(branch, "docs/") || strings.Contains(branch, "doc/") {
		return "books"
	}
	if strings.Contains(branch, "refactor/") {
		return "eyes"
	}
	if strings.Contains(branch, "test/") {
		return "rotating_light"
	}

	// デフォルト
	return "construction"
}

// generateSmartCommitMessage スマートなコミットメッセージを生成する
func generateSmartCommitMessage(hint string, files []string, branch string) string {
	if hint != "" {
		return hint
	}

	// ファイル拡張子から推測
	hasGo := false
	hasDoc := false
	hasTest := false

	for _, file := range files {
		if strings.HasSuffix(file, ".go") {
			hasGo = true
		}
		if strings.HasSuffix(file, ".md") || strings.HasSuffix(file, ".txt") {
			hasDoc = true
		}
		if strings.Contains(file, "test") || strings.HasSuffix(file, "_test.go") {
			hasTest = true
		}
	}

	if hasTest {
		return "Add or update tests"
	}
	if hasDoc {
		return "Update documentation"
	}
	if hasGo {
		return "Update Go implementation"
	}

	return "Update files"
}

// handleGitStatus Gitステータスの確認を処理する
func handleGitStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 変更されたファイルを取得
	changedFiles, err := git.GetChangedFilesList()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get changed files: %v", err)), nil
	}

	// ステージされたファイルを取得
	stagedFiles, err := git.GetStagedFiles()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get staged files: %v", err)), nil
	}

	// ブランチ情報を取得
	branch, err := git.GetBranch()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get branch: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Git Status:\nBranch: %s\nChanged files: %v\nStaged files: %v", branch, changedFiles, stagedFiles)), nil
}

// handleGitCommit 絵文字付きコミットを処理する
func handleGitCommit(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	emoji, err := request.RequireString("emoji")
	if err != nil {
		return mcp.NewToolResultError("emoji parameter is required"), nil
	}

	message, err := request.RequireString("message")
	if err != nil {
		return mcp.NewToolResultError("message parameter is required"), nil
	}

	// 引数の取得
	args := request.GetArguments()

	// offline パラメータの取得（オプション）
	offline := false
	if offlineParam, exists := args["offline"]; exists {
		if offlineBool, ok := offlineParam.(bool); ok {
			offline = offlineBool
		}
	}

	// 既存のビジネスロジックを呼び出してコミットを実行
	cm := git.CommitMessage{
		Emoji:   emoji,
		Message: message,
	}

	err = git.CommitWithOfflineMode(cm, offline)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to commit: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully committed with emoji '%s' and message '%s'", emoji, message)), nil
}

// handleGitAddFiles ファイルのステージングを処理する
func handleGitAddFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 引数の取得
	args := request.GetArguments()

	// files パラメータの取得（オプション）
	var filesToAdd []string
	if filesParam, exists := args["files"]; exists {
		if filesArray, ok := filesParam.([]interface{}); ok {
			for _, file := range filesArray {
				if fileStr, ok := file.(string); ok {
					filesToAdd = append(filesToAdd, fileStr)
				}
			}
		}
	}

	// ファイルが指定されていない場合は全ての変更ファイルを追加
	if len(filesToAdd) == 0 {
		changedFiles, err := git.GetChangedFilesList()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get changed files: %v", err)), nil
		}
		filesToAdd = changedFiles
	}

	// ファイルをステージング
	err := git.AddFiles(filesToAdd)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to add files: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully staged files: %v", filesToAdd)), nil
}

// handleGetEditedFiles 変更されたファイル一覧を取得する
func handleGetEditedFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	changedFiles, err := git.GetChangedFilesList()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get edited files: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Edited files: %v", changedFiles)), nil
}

// handleGetCurrentBranch 現在のブランチ名を取得する
func handleGetCurrentBranch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	branch, err := git.GetBranch()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get current branch: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Current branch: %s", branch)), nil
}
