package mcp

import (
	"fmt"

	"github.com/HidemaruOwO/pummit/internal/variable"
	"github.com/mark3labs/mcp-go/server"
)

func StartServer() error {
	s := server.NewMCPServer(
		"pummit",
		variable.VERSION,
		server.WithRecovery(),
	)

	if err := registerGitTools(s); err != nil {
		return fmt.Errorf("register git tools: %w", err)
	}
	if err := registerConfigTools(s); err != nil {
		return fmt.Errorf("register config tools: %w", err)
	}
	if err := registerAliasTools(s); err != nil {
		return fmt.Errorf("register alias tools: %w", err)
	}
	if err := registerDoctorTools(s); err != nil {
		return fmt.Errorf("register doctor tools: %w", err)
	}

	return server.ServeStdio(s)
}
