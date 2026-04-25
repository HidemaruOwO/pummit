package mcp

import (
	"bytes"
	"context"
	"fmt"

	"github.com/BurntSushi/toml"
	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerConfigTools(s *server.MCPServer) error {
	tool := mcp.NewTool("config.get",
		mcp.WithDescription("Get configuration values from the pummit configuration."),
		mcp.WithString("key", mcp.Description("Configuration key to retrieve (optional)")),
	)
	s.AddTool(tool, handleConfigGet)
	return nil
}

func handleConfigGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	key, _ := args["key"].(string)
	service := usecase.NewConfigService(rootconfig.NewStore(""))
	text, err := configGetText(service, key)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(text), nil
}

func configGetText(service *usecase.ConfigService, key string) (string, error) {
	if key == "" {
		cfg, err := service.List()
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	value, err := service.Get(key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s: %v", key, value), nil
}
