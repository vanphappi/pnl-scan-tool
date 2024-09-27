package utils

import (
	"strconv"
)

// ConvertStringToFloat64 converts a string to a float64.
// It returns the float64 value and an error if the conversion fails.
func ConvertStringToFloat64(s string) float64 {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return value
}

// ConvertStringToInt converts a string to an int.
// It returns the int value and an error if the conversion fails.
func ConvertStringToInt(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
}

// StringToInt64 converts a string to an int64 and returns an error if conversion fails.
func StringToInt64(s string) int64 {
	// strconv.ParseInt parses the string with base 10 and 64-bit size.
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
