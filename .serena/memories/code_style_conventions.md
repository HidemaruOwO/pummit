# Pummit Code Style and Conventions

## Language Standards
- **Go version**: 1.23.0+ required
- **Code formatting**: Standard `gofmt` formatting applied
- **File naming**: `snake_case.go` format consistently used
- **Package naming**: Lower case, single words preferred

## Coding Conventions

### Indentation and Formatting
- **Indentation**: 2 spaces (as specified in AGENTS.md)
- **Maximum line length**: 80 characters
- **Import organization**: External modules first, then internal paths, alphabetized within each group

### Naming Conventions
- **Variables and functions**: camelCase (`getUserConfig`, `parseEmojiName`)
- **Types and structs**: PascalCase (`Config`, `CommitMessage`)
- **Constants**: UPPER_SNAKE_CASE (`DEFAULT_CONFIG_PATH`, `VERSION`)
- **Package names**: lowercase, descriptive (`config`, `alias`, `doctor`)

### Error Handling Patterns
- Handle errors immediately where they occur
- Wrap external errors with context using `fmt.Errorf`
- Use structured error types for different error categories
- Log errors before returning them to callers

```go
// Example error handling pattern
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### Comment Guidelines
- **Write all comments in English** (as specified in project rules)
- Focus on **why** the code exists, not what it does
- Include business logic explanations and external constraints
- Avoid obvious descriptions of implementation details

#### Required Comments
- Why the code was written (background, rationale)
- Business logic and rule explanations
- External dependencies and constraints
- Important notes for future developers
- Performance and security considerations

#### Prohibited Comments
- Comments that only describe what the code does
- Simple implementation details
- Information obvious from variable/function names
- Outdated or incorrect information

## Architecture Patterns

### Layered Architecture
Following the 4-layer pattern defined in `docs/architecture-v3.md`:
1. **User Interface Layer**: CLI commands and interactive prompts
2. **Application Layer**: Business logic coordination
3. **Business Logic Layer**: Core functionality (Git ops, config, aliases)
4. **Infrastructure Layer**: File system, network, logging

### Package Organization
```
internal/
├── cli/          # Command-line interface (Cobra commands)
├── config/       # Configuration management
├── git/          # Git operations
├── alias/        # Alias management
├── doctor/       # Diagnostics
├── mcp/          # MCP server functionality
├── prompt/       # Interactive UI components
└── utils/        # Shared utilities

pkg/
├── gitmoji/      # External API integration
└── logger/       # Logging utilities
```

### Function and Method Design
- Keep functions focused on single responsibilities
- Use descriptive parameter and return value names
- Prefer returning errors over panicking
- Use context for cancellation and timeouts where appropriate

## Testing Conventions
- Test file naming: `*_test.go`
- Test function naming: `TestFunctionName`
- Benchmark function naming: `BenchmarkFunctionName`
- Use table-driven tests for multiple test cases
- Target 80% test coverage for critical components

## Configuration and Settings
- Use TOML format for configuration files
- Store configuration in `~/.config/pummit/` (Unix) or `%APPDATA%\pummit` (Windows)
- Support environment variable overrides where appropriate
- Provide sensible defaults for all configuration options

## Security Best Practices
- Validate all user inputs
- Use `filepath.Join()` for cross-platform path construction
- Never log sensitive information
- Implement proper file permissions (0644 for config files)
- Use HTTPS for all external API communications

## Git Integration Patterns
- Use `os/exec` for Git command execution
- Implement proper command injection prevention
- Handle Git errors gracefully with user-friendly messages
- Support both online and offline operation modes

## UI/UX Consistency
- Use consistent color coding for status messages
- Implement proper terminal width handling
- Support both interactive and non-interactive modes
- Provide clear help text for all commands
- Use emoji consistently for visual feedback

## Internationalization Readiness
- Keep all user-facing strings in separate constants
- Use structured logging for debug information
- Design UI components to handle variable text lengths
- Consider right-to-left language support in future versions