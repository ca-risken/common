package strings

// TruncateString truncates the input string to the specified character length
// UTF-8 safe - counts Unicode code points, not bytes
// suffix: appended when truncated (e.g., "..." or empty string)
func TruncateString(input string, maxLength int, suffix string) string {
	runes := []rune(input)
	if len(runes) <= maxLength {
		return input
	}
	return string(runes[:maxLength]) + suffix
}
