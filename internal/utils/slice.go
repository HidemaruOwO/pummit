package utils

// ContainsString は文字列スライスに特定の文字列が含まれているかどうかを返す
// gitmojiが含まれているかどうかを判定するために使用
func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// RemoveString は文字列スライスから特定の文字列を削除します
// gitmojiを削除するために使用
func RemoveString(slice []string, s string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

// FilterString は条件に合致する文字列のみを残します
// gitmojiを見つけるために使用
func FilterString(slice []string, condition func(string) bool) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if condition(item) {
			result = append(result, item)
		}
	}
	return result
}
