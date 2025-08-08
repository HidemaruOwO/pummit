# Pummit Codebase Structure and Architecture

## Directory Structure
```
pummit/
├── main.go                    # Application entry point
├── internal/                  # Private application code
│   ├── cli/                  # CLI commands (Cobra-based)
│   │   ├── root.go          # Root command and main CLI setup
│   │   ├── mcp.go           # MCP server command
│   │   ├── doctor.go        # Diagnostics command
│   │   ├── migrate.go       # Configuration migration
│   │   ├── version.go       # Version command
│   │   └── alias/           # Alias management commands
│   │       ├── add.go       # Add alias command
│   │       ├── delete.go    # Delete alias command
│   │       ├── list.go      # List aliases command
│   │       └── reset.go     # Reset aliases command
│   ├── config/              # Configuration management
│   │   ├── config.go        # Main configuration logic
│   │   ├── toml.go         # TOML format support
│   │   ├── migration.go    # JSON to TOML migration
│   │   └── compatibility.go # Legacy JSON support
│   ├── git/                 # Git operations
│   │   ├── git.go          # Core Git functionality
│   │   └── git_test.go     # Git tests
│   ├── alias/               # Alias management
│   │   └── alias.go        # Alias logic and storage
│   ├── doctor/              # System diagnostics
│   │   ├── checker.go      # Diagnostic checks
│   │   └── system.go       # System information
│   ├── mcp/                 # MCP server functionality
│   │   ├── server.go       # MCP server setup and tool registration
│   │   ├── tool_git.go     # Git-related MCP tools
│   │   ├── tool_alias.go   # Alias-related MCP tools
│   │   ├── tool_config.go  # Config-related MCP tools
│   │   └── tool_doctor.go  # Diagnostics MCP tool
│   ├── prompt/              # Interactive UI components
│   │   ├── prompt.go       # Basic prompts and confirmations
│   │   └── editor_selector.go # Editor selection UI (Bubble Tea)
│   ├── emojis/              # Emoji processing
│   │   └── emojis.go       # Emoji resolution and mapping
│   ├── utils/               # Utility functions
│   │   └── slice.go        # Slice manipulation utilities
│   └── variable/            # Constants and embedded data
│       ├── consts.go       # Application constants
│       ├── config.json     # Default configuration template
│       └── discord-emojis.flat.json # Embedded emoji data
├── pkg/                     # Public packages
│   ├── gitmoji/            # Gitmoji API integration
│   │   └── gitmoji.go      # API client for gitmoji service
│   └── logger/             # Logging utilities
│       └── logger.go       # Logger implementation
└── docs/                   # Documentation
    ├── architecture-v3.md  # Comprehensive architecture specification
    ├── llm-context.md     # Project context for AI development
    └── README.md          # User documentation
```

## Layered Architecture Implementation

### 1. User Interface Layer
**Location**: `internal/cli/`, `internal/prompt/`
**Responsibilities**:
- Command-line argument parsing (Cobra framework)
- Interactive user prompts (Bubble Tea framework)
- Output formatting and display
- Error message presentation

### 2. Application Layer
**Location**: `internal/cli/root.go`, command handlers
**Responsibilities**:
- Coordinating business logic calls
- Managing application flow
- Handling command-level error processing
- Configuration initialization

### 3. Business Logic Layer
**Location**: `internal/git/`, `internal/alias/`, `internal/config/`, `internal/doctor/`
**Responsibilities**:
- Core Git operations (commit creation, file staging)
- Alias management (create, delete, resolve)
- Configuration management (load, save, migrate)
- System diagnostics and health checks

### 4. Infrastructure Layer
**Location**: `pkg/gitmoji/`, `pkg/logger/`, standard library usage
**Responsibilities**:
- External API communication (Gitmoji service)
- File system operations
- Network requests and responses
- Logging and monitoring

## Key Components and Data Flow

### Configuration System
1. **Initialization**: `config.Init()` called at application startup
2. **Path Resolution**: Determines config file location based on OS
3. **Format Detection**: Checks for TOML first, then JSON (legacy)
4. **Migration**: Automatic JSON to TOML conversion when needed
5. **Loading**: Deserializes configuration into `Config` struct
6. **Usage**: Accessed globally via `config.CurrentTOMLConfig`

### Git Operations
1. **File Detection**: `git diff --name-only --cached` to find staged files
2. **Emoji Resolution**: Input emoji/alias resolved via `alias.GetEmoji()`
3. **Message Formation**: Emoji + message + file list formatting
4. **Commit Execution**: `git commit -m "..."` via `os/exec`

### MCP Server Integration
1. **Server Setup**: JSON-RPC over stdio communication
2. **Tool Registration**: 9 tools registered for AI interaction
3. **Business Logic Reuse**: Direct calls to existing internal packages
4. **Smart Workflow**: `git.smart_commit` provides autonomous Git operations

### Alias System
1. **Storage**: Aliases stored in configuration file as array of arrays
2. **Resolution**: Multi-step lookup (exact match → alias → fallback)
3. **Management**: CRUD operations via dedicated alias commands
4. **Integration**: Used by both CLI and MCP interfaces

## Module Dependencies and Relationships

### Internal Dependencies
- `cli/` → `config/`, `git/`, `alias/`, `doctor/`
- `mcp/` → `git/`, `alias/`, `config/`, `doctor/`
- `git/` → `config/`, `alias/`, `emojis/`
- `config/` → `variable/` (for defaults)

### External Dependencies
- `cobra` for CLI framework
- `bubbletea` for interactive UI
- `toml` for configuration parsing
- `mcp-go` for AI integration protocol

## Data Structures

### Core Types
```go
type Config struct {
    WriteEmoji            bool          // Output emoji directly vs :name:
    UseAlias             bool          // Enable alias resolution
    UseLimitPathesLength bool          // Limit file path display
    LimitPathesLength    int           // Max characters for file paths
    Alias                [][]string    // Alias definitions
}

type CommitMessage struct {
    Emoji   string    // Resolved emoji for commit
    Message string    // User-provided message
}
```

### Configuration File Formats
**TOML (Current)**:
```toml
[base]
emoji = true
filesLength = 50

[[alias.entries]]
shortcuts = ["s", "feat", "feature"]
name = "sparkles"
emoji = "✨"
```

**JSON (Legacy)**:
```json
{
  "writeEmoji": true,
  "alias": [["s,feat,feature", "sparkles", "✨"]]
}
```

This structure provides a clear separation of concerns while maintaining efficient communication between layers and components.