# Pummit Tech Stack and Dependencies

## Programming Language
- **Go**: 1.23.0+ required
- **Target platforms**: Windows, macOS, Linux (cross-platform)

## Major Dependencies

### CLI and User Interface
- **github.com/spf13/cobra v1.9.1**: CLI framework for command structure
- **github.com/charmbracelet/bubbletea v1.3.6**: TUI framework for interactive mode
- **github.com/charmbracelet/bubbles v0.20.0**: TUI components (list selection, etc.)
- **github.com/charmbracelet/lipgloss v1.1.0**: TUI styling and formatting
- **github.com/fatih/color v1.18.0**: Terminal color output

### Configuration and Data
- **github.com/BurntSushi/toml v1.5.0**: TOML configuration file handling
- **github.com/mark3labs/mcp-go v0.37.0**: MCP (Model Context Protocol) server implementation

### Standard Library Usage
- **os/exec**: Git command execution
- **encoding/json**: JSON configuration migration
- **net/http**: Gitmoji API communication
- **path/filepath**: Cross-platform path handling

## External Services
- **Gitmoji API**: `https://raw.githubusercontent.com/carloscuesta/gitmoji/master/packages/gitmojis/src/gitmojis.json`
- **Fallback**: Embedded Discord emoji data for offline operation

## Build System
- **Go modules**: Standard dependency management
- **GoReleaser**: Automated release building (`.goreleaser.yaml`)
- **GitHub Actions**: CI/CD pipeline (`.github/workflows/`)

## Development Tools
- **golangci-lint**: Code linting (referenced in AGENTS.md)
- **go test**: Unit testing framework
- **go build**: Standard build process

## Architecture Patterns
- **Layered architecture**: UI/Application/Business/Infrastructure
- **Cobra CLI pattern**: Command and subcommand structure
- **MCP protocol**: JSON-RPC over stdio for AI integration

## Security Considerations
- HTTPS-only for external API calls
- Path validation for configuration files
- Input sanitization for commit messages
- Minimal privilege principle

## Performance Targets
- **Startup time**: < 100ms for basic operations
- **Memory usage**: < 10MB for standard operations
- **Network timeout**: 3 seconds for API calls with fallback