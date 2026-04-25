package variable

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestVersionMatchesModuleMajor(t *testing.T) {
	root := filepath.Join("..", "..")
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}

	moduleLine := ""
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			moduleLine = strings.TrimPrefix(line, "module ")
			break
		}
	}
	if moduleLine == "" {
		t.Fatal("module line not found in go.mod")
	}

	major := versionMajor(t, VERSION)
	moduleMajor := modulePathMajor(moduleLine)

	if major == 1 && moduleMajor != 1 {
		t.Fatalf("VERSION=%q implies module major v1, but go.mod module path is %q", VERSION, moduleLine)
	}

	if major >= 2 && moduleMajor != major {
		t.Fatalf("VERSION=%q requires module path major /v%d, but go.mod module path is %q", VERSION, major, moduleLine)
	}

	if moduleMajor >= 2 && major != moduleMajor {
		t.Fatalf("go.mod module path %q requires VERSION major v%d, but VERSION=%q", moduleLine, moduleMajor, VERSION)
	}
}

func versionMajor(t *testing.T, version string) int {
	t.Helper()

	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		t.Fatalf("VERSION=%q is not semver-like", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		t.Fatalf("parse VERSION major %q: %v", version, err)
	}

	return major
}

func modulePathMajor(modulePath string) int {
	re := regexp.MustCompile(`/v([0-9]+)$`)
	matches := re.FindStringSubmatch(modulePath)
	if len(matches) != 2 {
		return 1
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return 1
	}

	return major
}
