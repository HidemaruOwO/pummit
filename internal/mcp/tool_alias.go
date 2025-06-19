package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/HidemaruOwO/pummit/internal/alias"
)

// registerAliasTools エイリアス関連のツールをMCPサーバーに登録する
func registerAliasTools(s *server.MCPServer) error {
	// alias.list ツール: エイリアス一覧取得
	listTool := mcp.NewTool("alias.list",
		mcp.WithDescription("Get a list of all available emoji aliases and their corresponding emojis."),
	)
	s.AddTool(listTool, handleAliasList)

	return nil
}

// handleAliasList エイリアス一覧を取得する
func handleAliasList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 既存のビジネスロジックを使用してエイリアス一覧を取得
	aliases := alias.List()

	return mcp.NewToolResultText(fmt.Sprintf("Available aliases: %+v", aliases)), nil
}
