package git

import "testing"

// マルチバイト文字を含む文字列の切り詰め処理を検証
func TestTruncateWithEllipsis(t *testing.T) {
	t.Run("ascii", func(t *testing.T) {
		s := "abcde"
		got := truncateWithEllipsis(s, 3)
		want := "abc..."
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	})

	t.Run("multibyte", func(t *testing.T) {
		s := "あいうえお"
		got := truncateWithEllipsis(s, 3)
		want := "あいう..."
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	})

	t.Run("no trunc", func(t *testing.T) {
		s := "abc"
		got := truncateWithEllipsis(s, 5)
		if got != s {
			t.Fatalf("expected original string, got %q", got)
		}
	})
}
