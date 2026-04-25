package commit

import "testing"

func TestFormatMessage(t *testing.T) {
	got := FormatMessage("✨", "Add x", "a.go, b.go")
	want := "✨ Add x (a.go, b.go)"
	if got != want {
		t.Fatalf("FormatMessage() = %q, want %q", got, want)
	}
}

func TestTruncateFiles(t *testing.T) {
	tests := []struct {
		name  string
		files string
		limit int
		want  string
	}{
		{name: "no truncation", files: "a.go", limit: 10, want: "a.go"},
		{name: "zero limit", files: "a.go, b.go", limit: 0, want: "a.go, b.go"},
		{name: "truncate ascii", files: "a.go, b.go", limit: 4, want: "a.go..."},
		{name: "truncate multibyte", files: "あいうえお", limit: 3, want: "あいう..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateFiles(tt.files, tt.limit)
			if got != tt.want {
				t.Fatalf("TruncateFiles(%q, %d) = %q, want %q", tt.files, tt.limit, got, tt.want)
			}
		})
	}
}
