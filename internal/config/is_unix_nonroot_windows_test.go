//go:build windows

package config

// isUnixNonRoot always returns false on Windows.
func isUnixNonRoot() bool { return false }
