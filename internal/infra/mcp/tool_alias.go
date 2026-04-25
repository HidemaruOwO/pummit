package mcp

import (
	"context"
	"fmt"
	"strings"

	domainconfig "github.com/HidemaruOwO/pummit/internal/domain/config"
	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	rootemoji "github.com/HidemaruOwO/pummit/internal/infra/emoji"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerAliasTools(s *server.MCPServer) error {
	tool := mcp.NewTool("alias.list",
		mcp.WithDescription("Get a list of all available emoji aliases and their corresponding emojis."),
	)
	s.AddTool(tool, handleAliasList)
	return nil
}

func handleAliasList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	configService := usecase.NewConfigService(rootconfig.NewStore(""))
	service := usecase.NewAliasService(configService, rootemoji.NewCatalog(), rootemoji.NewGitmojiClient())
	entries, err := service.List()
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(formatAliasEntries(entries)), nil
}

func formatAliasEntries(entries []domainconfig.AliasEntry) string {
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, fmt.Sprintf("%s %s => %s", entry.Emoji, entry.Name, strings.Join(entry.Shortcuts, ",")))
	}
	return strings.Join(lines, "\n")
}
