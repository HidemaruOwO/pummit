package cli

import (
	rootmcp "github.com/HidemaruOwO/pummit/internal/infra/mcp"
	"github.com/spf13/cobra"
)

func newMCPCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Start the MCP (Model Context Protocol) server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootmcp.StartServer()
		},
	}
}
