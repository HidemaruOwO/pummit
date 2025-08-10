# Pummit Development Commands and Scripts

## Build Commands
```bash
# Standard build
go build -o pummit .

# Cross-platform build with optimization
go build -ldflags="-s -w" -trimpath -o pummit .

# Build to specific output directory
go build -o /tmp/pummit .
```

## Testing Commands
```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

## Linting and Code Quality
```bash
# Run golangci-lint (as specified in AGENTS.md)
golangci-lint run

# Format code
go fmt ./...

# Vet code
go vet ./...

# Check for vulnerabilities
govulncheck ./...
```

## Development Workflow
```bash
# Install dependencies
go mod download

# Update dependencies
go mod tidy

# Vendor dependencies (if needed)
go mod vendor
```

## Running the Application
```bash
# Basic commit
./pummit sparkles "Add new feature"

# With offline mode
./pummit --offline bug "Fix authentication issue"

# Interactive mode (planned)
./pummit interactive

# MCP server mode
./pummit mcp

# Diagnostics
./pummit doctor

# Migration
./pummit migrate
```

## Release Process
```bash
# Using GoReleaser (automated)
goreleaser release --clean

# Manual release build
goreleaser build --single-target --clean
```

## Configuration Management
```bash
# View current configuration
cat ~/.config/pummit/config.toml

# Edit configuration (planned unified command)
pummit config

# Migrate from JSON to TOML
pummit migrate
```

## Debugging and Diagnostics
```bash
# System diagnostics
pummit doctor

# Version information
pummit --version

# Help information
pummit --help
pummit [command] --help
```

## Docker/Container Support
Currently not implemented, but can be added:
```dockerfile
# Potential Dockerfile structure
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o pummit .

FROM alpine:latest
RUN apk --no-cache add git ca-certificates
WORKDIR /root/
COPY --from=builder /app/pummit .
CMD ["./pummit"]
```

## macOS Specific
```bash
# Homebrew installation (planned)
brew tap hidemaruowo/tap
brew install pummit
```

## Windows Specific
```powershell
# PowerShell execution
.\pummit.exe sparkles "Add new feature"

# Using Scoop (planned)
scoop install pummit
```

## CI/CD Integration
The project uses GitHub Actions for:
- Automated testing on multiple platforms
- Build verification
- Release automation
- Dependency updates via Dependabot