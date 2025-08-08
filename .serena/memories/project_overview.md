# Pummit Project Overview

## Project Purpose
Pummit is a professional Git commit message CLI tool that helps developers create beautiful, consistent, and meaningful commit messages with emoji support and smart automation.

## Core Features (v2.0.0 Implemented)
- **Emoji-prefixed commits**: Visual commit messages with emoji prefixes
- **Smart automation**: Branch-based emoji suggestions (feature/ → ✨, fix/ → 🐛)
- **Offline support**: `--offline` flag for network-independent operation
- **Flexible configuration**: TOML-based configuration with automatic migration from JSON
- **Alias system**: Multiple shortcuts for emojis (e.g., "s,feat,feature" → "✨")
- **Comprehensive diagnostics**: `pummit doctor` command for environment checking
- **AI integration**: MCP server functionality for Claude Desktop integration
- **Interactive UI**: Bubble Tea-based confirmation dialogs and prompts

## Target Users
- Developers who want consistent, visually appealing commit messages
- Teams following Conventional Commits or Gitmoji standards
- Users who prefer CLI tools with smart automation
- Organizations using AI-assisted development workflows

## Key Benefits
- Improved commit message consistency across teams
- Visual clarity with emoji prefixes
- Smart branch-based suggestions reduce manual work
- AI integration enables natural language Git operations
- Cross-platform compatibility (Windows, macOS, Linux)

## Technology Stack
- **Language**: Go 1.23.0+
- **CLI Framework**: Cobra
- **TUI Framework**: Bubble Tea
- **Configuration**: TOML (BurntSushi/toml)
- **AI Integration**: MCP (Model Context Protocol)
- **External APIs**: Gitmoji API (with offline fallback)

## Project Status
Currently in v2.0.0 with major features implemented. Planning v3.0.0 with CLI structure unification and comprehensive test coverage.