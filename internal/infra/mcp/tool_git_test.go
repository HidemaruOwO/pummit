package mcp

import "testing"

func TestGenerateSmartCommitMessage(t *testing.T) {
	tests := []struct {
		name  string
		hint  string
		files []string
		want  string
	}{
		{name: "hint wins", hint: "Custom message", files: []string{"file.go"}, want: "Custom message"},
		{name: "test files", files: []string{"service_test.go"}, want: "Add or update tests"},
		{name: "docs files", files: []string{"README.md"}, want: "Update documentation"},
		{name: "go files", files: []string{"main.go"}, want: "Update Go implementation"},
		{name: "fallback", files: []string{"package.json"}, want: "Update files"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateSmartCommitMessage(tt.hint, tt.files)
			if got != tt.want {
				t.Fatalf("generateSmartCommitMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatGitStatus(t *testing.T) {
	got := formatGitStatus(gitStatus{Branch: "main", Changed: []string{"a.go"}, Staged: []string{"b.go"}})
	want := "Branch: main\nChanged files: [a.go]\nStaged files: [b.go]"
	if got != want {
		t.Fatalf("formatGitStatus() = %q, want %q", got, want)
	}
}
