package utils

import (
	"regexp"
	"strings"
)

// PreprocessJSON takes a string that is supposed to be valid JSON, but might have some
// issues (like single quotes instead of double quotes, or missing quotes around keys).
// It fixes these issues and returns the fixed string.
func PreprocessJSON(input string) string {
	// Replace single quotes with double quotes
	input = strings.ReplaceAll(input, "'", "\"")

	// Handle boolean values
	input = strings.ReplaceAll(input, ": true", ": true")
	input = strings.ReplaceAll(input, ": false", ": false")

	// Handle null values
	input = strings.ReplaceAll(input, ": null", ": null")

	// Handle undefined values
	input = strings.ReplaceAll(input, ": undefined", ": null")

	// Remove trailing commas in objects and arrays
	input = regexp.MustCompile(`,\s*([}\]])`).ReplaceAllString(input, "$1")

	// Wrap unquoted keys in double quotes
	input = regexp.MustCompile(`(\w+):`).ReplaceAllString(input, "\"$1\":")

	// Handle function values
	input = regexp.MustCompile(`:\s*function\s*\([^)]*\)\s*{[^}]*}`).ReplaceAllString(input, ": \"[Function]\"")

	return input
}
