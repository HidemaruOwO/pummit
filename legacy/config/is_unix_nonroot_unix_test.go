//go:build !windows

package config

import (
	"os"
)

// isUnixNonRoot reports whether tests run on a Unix-like system without root privileges.
func isUnixNonRoot() bool {
	return os.Geteuid() != 0
}
