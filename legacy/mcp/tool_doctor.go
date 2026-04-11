package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/HidemaruOwO/pummit/legacy/doctor"
)

// registerDoctorTools 診断関連のツールをMCPサーバーに登録する
func registerDoctorTools(s *server.MCPServer) error {
	// doctor.check ツール: システム診断
	checkTool := mcp.NewTool("doctor.check",
		mcp.WithDescription("Run comprehensive system diagnostics to check Git configuration, repository status, config files, network connectivity, and file permissions."),
	)
	s.AddTool(checkTool, handleDoctorCheck)

	return nil
}

// handleDoctorCheck システム診断を実行する
func handleDoctorCheck(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 既存のビジネスロジックを使用してシステム診断を実行し、フォーマット済みの結果を取得
	output := doctor.FormatDiagnosticResults()

	return mcp.NewToolResultText(output), nil
}
