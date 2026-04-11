package utils

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}

func FilterString(slice []string, condition func(string) bool) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if condition(item) {
			result = append(result, item)
		}
	}
	return result
}
