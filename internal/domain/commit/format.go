package commit

import "fmt"

func FormatMessage(prefix, message, files string) string {
	return fmt.Sprintf("%s %s (%s)", prefix, message, files)
}

func TruncateFiles(files string, limit int) string {
	if limit <= 0 {
		return files
	}

	runes := []rune(files)
	if len(runes) <= limit {
		return files
	}

	return string(runes[:limit]) + "..."
}
