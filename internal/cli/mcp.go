package cli

import (
	"github.com/spf13/cobra"

	"github.com/HidemaruOwO/pummit/internal/mcp"
	"github.com/HidemaruOwO/pummit/pkg/logger"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the MCP (Model Context Protocol) server",
	Long: `Start the MCP server for integration with AI assistants.

The MCP server provides a standardized way for AI assistants to interact with pummit's
Git operations, configuration management, and diagnostic tools through JSON-RPC over stdio.

Usage with Claude Desktop:
1. Add this configuration to your Claude Desktop MCP settings
2. The server will handle Git operations, alias management, and system diagnostics
3. AI assistants can then use natural language to perform Git commits and other operations

Available MCP tools:
- git.smart_commit: Intelligent commit workflow with auto-analysis
- git.status: Repository status information  
- git.commit: Direct commit with specified parameters
- git.add_files: Stage files for commit
- git.get_edited_files: List modified files
- git.get_current_branch: Get current branch name
- alias.list: List available emoji aliases
- config.get: Retrieve configuration values
- doctor.check: Run system diagnostics`,
	Run: runMcpCommand,
}

// runMcpCommand MCPサーバーを起動する
func runMcpCommand(cmd *cobra.Command, args []string) {
	log := logger.New()

	log.Info("Starting MCP server...")

	// MCPサーバーを起動（stdio モード）
	if err := mcp.StartMCPServer(); err != nil {
		log.Error("Failed to start MCP server: " + err.Error())
		return
	}
}

func init() {
	// ルートコマンドにmcpサブコマンドを追加
	rootCmd.AddCommand(mcpCmd)
}
