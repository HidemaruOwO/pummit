# Legacy Implementation

This directory contains the pre-rewrite implementation of `pummit`.

Rules:
- New feature work should target the new `internal/` layout, not `legacy/`
- `legacy/` remains as a reference and compatibility layer during the rewrite
- The root `main.go` may continue to call into legacy code until the new implementation fully replaces it
