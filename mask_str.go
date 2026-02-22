package std

import "strings"

// MaskStr
//
// Examples:
//   - "" --> ""
//   - "s" --> "*"
//   - "Secret-codE" --> "***********"
func MaskStr(s string) string {
	if s == "" {
		return s
	}

	return strings.Repeat("*", len([]rune(s)))
}

// MaskStrNotFirst
//
// Examples:
//   - "" --> ""
//   - "s" --> "s"
//   - "Secret-codE" --> "S**********"
func MaskStrNotFirst(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	lr := len(runes)

	return string(runes[0]) + strings.Repeat("*", lr-1)
}

// MaskStrNotFirstLast
//
// Examples:
//   - "" --> ""
//   - "s" --> "s"
//   - "se" --> "se"
//   - "Secret-codE" --> "S*********E"
func MaskStrNotFirstLast(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	lr := len(runes)

	if lr <= 2 {
		return s
	}

	return string(runes[0]) + strings.Repeat("*", lr-2) + string(runes[len(runes)-1])
}
