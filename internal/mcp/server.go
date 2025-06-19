package mcp

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/HidemaruOwO/pummit/internal/variable"
)

// StartMCPServer MCPサーバーを起動する
// JSON-RPC over stdio によるMCPプロトコル対応を提供する
func StartMCPServer() error {
	// MCPサーバーの初期化
	s := server.NewMCPServer(
		"pummit",
		variable.VERSION,      // 動的にバージョン情報を設定
		server.WithRecovery(), // ハンドラ内でのパニックから回復
	)

	// 各ツールをサーバーに登録
	if err := registerGitTools(s); err != nil {
		return fmt.Errorf("failed to register git tools: %w", err)
	}

	if err := registerConfigTools(s); err != nil {
		return fmt.Errorf("failed to register config tools: %w", err)
	}

	if err := registerAliasTools(s); err != nil {
		return fmt.Errorf("failed to register alias tools: %w", err)
	}

	if err := registerDoctorTools(s); err != nil {
		return fmt.Errorf("failed to register doctor tools: %w", err)
	}

	// デバッグログの設定
	if os.Getenv("PUMMIT_MCP_DEBUG") == "1" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("Starting MCP server (version: %s)", variable.VERSION)
	}

	// 標準入出力でサーバーを起動
	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("mcp server error: %w", err)
	}

	return nil
}

// NewToolResultJSON JSON形式の結果を作成するヘルパー関数
func NewToolResultJSON(data interface{}) *mcp.CallToolResult {
	return mcp.NewToolResultText(fmt.Sprintf("%+v", data))
}
