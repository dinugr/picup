package utils

func GetStrMaxLength(arr []string) int {
	maxLen := 0
	for _, s := range arr {
		if len(s) > maxLen {
			maxLen = len(s)
		}
	}
	return maxLen
}
