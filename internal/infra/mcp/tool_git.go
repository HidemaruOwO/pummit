package mcp

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	rootemoji "github.com/HidemaruOwO/pummit/internal/infra/emoji"
	rootgit "github.com/HidemaruOwO/pummit/internal/infra/git"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerGitTools(s *server.MCPServer) error {
	statusTool := mcp.NewTool("git.status",
		mcp.WithDescription("Get the current status of the Git repository, including staged and unstaged files."),
	)
	s.AddTool(statusTool, handleGitStatus)

	commitTool := mcp.NewTool("git.commit",
		mcp.WithDescription("Create a Git commit with explicit emoji and message."),
		mcp.WithString("emoji", mcp.Required(), mcp.Description("Emoji name or alias to use for the commit")),
		mcp.WithString("message", mcp.Required(), mcp.Description("Commit message")),
	)
	s.AddTool(commitTool, handleGitCommit)

	smartCommitTool := mcp.NewTool("git.smart_commit",
		mcp.WithDescription("Analyze repository changes and create a commit using branchMapping and generated message."),
		mcp.WithString("message_hint", mcp.Description("Optional hint for commit message")),
		mcp.WithBoolean("auto_approve", mcp.Description("Execute commit immediately (default: false)")),
	)
	s.AddTool(smartCommitTool, handleGitSmartCommit)

	return nil
}

func handleGitStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	status, err := collectGitStatus()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(formatGitStatus(status)), nil
}

func handleGitCommit(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	emojiName, err := request.RequireString("emoji")
	if err != nil {
		return mcp.NewToolResultError("emoji parameter is required"), nil
	}
	message, err := request.RequireString("message")
	if err != nil {
		return mcp.NewToolResultError("message parameter is required"), nil
	}

	result, err := commitWithExplicitArgs(emojiName, message)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

func handleGitSmartCommit(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	messageHint := ""
	if raw, ok := args["message_hint"].(string); ok {
		messageHint = raw
	}
	autoApprove := false
	if raw, ok := args["auto_approve"].(bool); ok {
		autoApprove = raw
	}

	result, err := smartCommit(messageHint, autoApprove)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(result), nil
}

type gitStatus struct {
	Branch  string
	Changed []string
	Staged  []string
}

func collectGitStatus() (gitStatus, error) {
	wd, err := os.Getwd()
	if err != nil {
		return gitStatus{}, err
	}
	client := rootgit.NewClient(wd)
	branch, err := client.CurrentBranch()
	if err != nil {
		return gitStatus{}, err
	}
	changed, err := client.ChangedFiles()
	if err != nil {
		return gitStatus{}, err
	}
	staged, err := client.StagedFiles()
	if err != nil {
		return gitStatus{}, err
	}
	sort.Strings(changed)
	sort.Strings(staged)
	return gitStatus{Branch: branch, Changed: changed, Staged: staged}, nil
}

func formatGitStatus(status gitStatus) string {
	return fmt.Sprintf("Branch: %s\nChanged files: %v\nStaged files: %v", status.Branch, status.Changed, status.Staged)
}

func commitWithExplicitArgs(emojiName, message string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	service := usecase.NewCommitService(
		usecase.NewConfigService(rootconfig.NewStore("")),
		rootgit.NewClient(wd),
		rootemoji.NewCatalog(),
		rootemoji.NewGitmojiClient(),
	)
	result, err := service.CommitExplicit(emojiName, message)
	if err != nil {
		return "", err
	}
	if !result.Committed {
		return "No staged files to commit", nil
	}
	return result.Message, nil
}

func smartCommit(messageHint string, autoApprove bool) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	client := rootgit.NewClient(wd)
	branch, err := client.CurrentBranch()
	if err != nil {
		return "", err
	}
	changed, err := client.ChangedFiles()
	if err != nil {
		return "", err
	}
	staged, err := client.StagedFiles()
	if err != nil {
		return "", err
	}
	if len(staged) == 0 && len(changed) > 0 {
		if err := client.AddFiles(changed); err != nil {
			return "", err
		}
		staged, err = client.StagedFiles()
		if err != nil {
			return "", err
		}
	}
	if len(staged) == 0 {
		return "No changes to commit. Repository is clean.", nil
	}
	message := generateSmartCommitMessage(messageHint, staged)
	if !autoApprove {
		return fmt.Sprintf("Commit plan ready:\nBranch: %s\nStaged files: %v\nSuggested message: %s\n\nCall git.smart_commit with auto_approve=true to execute.", branch, staged, message), nil
	}

	service := usecase.NewCommitService(
		usecase.NewConfigService(rootconfig.NewStore("")),
		client,
		rootemoji.NewCatalog(),
		rootemoji.NewGitmojiClient(),
	)
	result, err := service.CommitAuto(message)
	if err != nil {
		return "", err
	}
	if !result.Committed {
		return "No staged files to commit", nil
	}
	return result.Message, nil
}

func generateSmartCommitMessage(hint string, files []string) string {
	if hint != "" {
		return hint
	}
	hasDoc := false
	hasTest := false
	hasGo := false
	for _, file := range files {
		switch {
		case strings.HasSuffix(file, "_test.go") || strings.Contains(file, "test"):
			hasTest = true
		case strings.HasSuffix(file, ".md") || strings.HasSuffix(file, ".txt"):
			hasDoc = true
		case strings.HasSuffix(file, ".go"):
			hasGo = true
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
