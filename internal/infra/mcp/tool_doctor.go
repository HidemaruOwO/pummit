package mcp

import (
	"context"
	"fmt"
	"strings"

	rootconfig "github.com/HidemaruOwO/pummit/internal/infra/config"
	"github.com/HidemaruOwO/pummit/internal/usecase"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerDoctorTools(s *server.MCPServer) error {
	tool := mcp.NewTool("doctor.check",
		mcp.WithDescription("Run diagnostics for Git, config, repository, network, and environment."),
	)
	s.AddTool(tool, handleDoctorCheck)
	return nil
}

func handleDoctorCheck(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	service := usecase.NewDoctorService(usecase.NewConfigService(rootconfig.NewStore("")), rootconfig.NewStore(""))
	report := service.Run()
	return mcp.NewToolResultText(formatDoctorReport(report)), nil
}

func formatDoctorReport(report usecase.DoctorReport) string {
	lines := make([]string, 0, len(report.System)+len(report.Checks))
	for _, item := range report.System {
		lines = append(lines, fmt.Sprintf("%s: %s", item.Name, item.Message))
	}
	for _, item := range report.Checks {
		lines = append(lines, fmt.Sprintf("%s [%s]: %s", item.Name, item.Status, item.Message))
	}
	return strings.Join(lines, "\n")
}
