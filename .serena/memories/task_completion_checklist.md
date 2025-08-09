# Pummit Task Completion Checklist

## When Completing Any Development Task

### 1. Code Quality Verification
- [ ] Run `golangci-lint run` to ensure code quality standards
- [ ] Execute `go test ./...` to verify all tests pass
- [ ] Check `go fmt ./...` for proper code formatting
- [ ] Run `go vet ./...` for static analysis

### 2. Build Verification
- [ ] Execute `go build -o pummit .` to ensure clean compilation
- [ ] Test the built binary with basic operations
- [ ] Verify cross-platform compatibility if changes affect OS-specific code

### 3. Feature-Specific Testing
- [ ] Test the new feature in isolation
- [ ] Verify integration with existing functionality
- [ ] Test error cases and edge conditions
- [ ] Ensure offline mode compatibility (if applicable)

### 4. Configuration and Migration
- [ ] Test configuration changes with both TOML and JSON formats (during transition period)
- [ ] Verify migration functionality if config changes are involved
- [ ] Test with default configuration
- [ ] Validate configuration file permissions and security

### 5. Documentation Updates
- [ ] Update relevant documentation in `docs/` directory
- [ ] Update `docs/llm-context.md` with new insights or architectural changes
- [ ] Update `README.md` if user-facing features change
- [ ] Update `AGENTS.md` if development commands change

### 6. AI Integration (MCP) Testing
- [ ] Test MCP server functionality with `pummit mcp` if Git operations are modified
- [ ] Verify that changes don't break natural language Git operations
- [ ] Test the `git.smart_commit` workflow if core Git functionality is affected

### 7. Backward Compatibility
- [ ] Ensure existing configuration files continue to work
- [ ] Test migration from previous versions
- [ ] Verify alias functionality remains intact
- [ ] Test offline emoji data fallback

### 8. Performance Verification
- [ ] Measure startup time (should remain < 100ms for basic operations)
- [ ] Check memory usage for any significant increases
- [ ] Verify network timeout behavior (3 seconds for API calls)

### 9. Security Review
- [ ] Review any new file operations for security implications
- [ ] Validate input sanitization for user-provided data
- [ ] Ensure no sensitive information is logged
- [ ] Check file permissions for created files

### 10. Platform-Specific Testing
- [ ] Test on Unix-like systems (macOS/Linux)
- [ ] Test on Windows (if changes affect path handling or commands)
- [ ] Verify terminal compatibility (especially for UI components)

## Before Committing Changes

### 1. Pre-commit Verification
- [ ] All tests passing: `go test ./... -race`
- [ ] Code properly formatted: `go fmt ./...`
- [ ] No linting issues: `golangci-lint run`
- [ ] Clean build: `go build -o pummit .`

### 2. Commit Message Standards
- [ ] Use pummit itself to create the commit message (dogfooding)
- [ ] Follow emoji prefix conventions
- [ ] Include clear, descriptive commit message
- [ ] Reference issues or features if applicable

### 3. Branch Management
- [ ] Use appropriate branch naming (`feature/`, `fix/`, `refactor/`)
- [ ] Ensure branch is up to date with main
- [ ] Squash commits if necessary for clean history

## Release Preparation (for maintainers)

### 1. Version Management
- [ ] Update version numbers in relevant files
- [ ] Update `internal/variable/consts.go` if version constant exists
- [ ] Update changelog or release notes

### 2. Comprehensive Testing
- [ ] Run full test suite on multiple platforms
- [ ] Test installation from scratch
- [ ] Verify GoReleaser configuration
- [ ] Test package managers (Homebrew, etc.) if applicable

### 3. Documentation Completeness
- [ ] All new features documented in README
- [ ] Architecture documentation updated if structural changes made
- [ ] API documentation updated for MCP changes
- [ ] Migration guides updated if breaking changes introduced

## Troubleshooting Common Issues

### Build Failures
- Check Go version compatibility (1.23.0+)
- Verify all dependencies are available: `go mod download`
- Clean module cache if persistent issues: `go clean -modcache`

### Test Failures
- Run tests with verbose output: `go test -v ./...`
- Check for race conditions: `go test -race ./...`
- Verify test environment setup (Git configuration, etc.)

### Linting Issues
- Fix formatting first: `go fmt ./...`
- Address vet warnings: `go vet ./...`
- Review golangci-lint specific errors and warnings

This checklist ensures consistent quality and reduces the risk of regressions when completing development tasks.