package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/HidemaruOwO/pummit/internal/config"
)

// registerConfigTools 設定関連のツールをMCPサーバーに登録する
func registerConfigTools(s *server.MCPServer) error {
	// config.get ツール: 設定値取得
	getTool := mcp.NewTool("config.get",
		mcp.WithDescription("Get configuration values from the pummit configuration."),
		mcp.WithString("key", mcp.Description("Configuration key to retrieve (optional, returns all if not specified)")),
	)
	s.AddTool(getTool, handleConfigGet)

	return nil
}

// handleConfigGet 設定値を取得する
func handleConfigGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 引数の取得
	args := request.GetArguments()

	// key パラメータの取得（オプション）
	key := ""
	if keyParam, exists := args["key"]; exists {
		if keyStr, ok := keyParam.(string); ok {
			key = keyStr
		}
	}

	// 現在の設定を取得
	currentConfig := config.CurrentTOMLConfig

	// 特定のキーが指定された場合
	if key != "" {
		switch key {
		case "base.emoji":
			return mcp.NewToolResultText(fmt.Sprintf("base.emoji: %t", currentConfig.Base.Emoji)), nil
		case "base.filesLength":
			return mcp.NewToolResultText(fmt.Sprintf("base.filesLength: %d", currentConfig.Base.FilesLength)), nil
		case "alias.enabled":
			return mcp.NewToolResultText(fmt.Sprintf("alias.enabled: %t", currentConfig.Alias.Enabled)), nil
		case "interactive.enabled":
			return mcp.NewToolResultText(fmt.Sprintf("interactive.enabled: %t", currentConfig.Interactive.Enabled)), nil
		case "locale.language":
			return mcp.NewToolResultText(fmt.Sprintf("locale.language: %s", currentConfig.Locale.Language)), nil
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unknown configuration key: %s", key)), nil
		}
	}

	// 全設定を返す
	configSummary := fmt.Sprintf(`Current Configuration:
- base.emoji: %t
- base.filesLength: %d
- alias.enabled: %t
- interactive.enabled: %t
- locale.language: %s
- Number of alias entries: %d`,
		currentConfig.Base.Emoji,
		currentConfig.Base.FilesLength,
		currentConfig.Alias.Enabled,
		currentConfig.Interactive.Enabled,
		currentConfig.Locale.Language,
		len(currentConfig.Alias.Entries))

	return mcp.NewToolResultText(configSummary), nil
}
