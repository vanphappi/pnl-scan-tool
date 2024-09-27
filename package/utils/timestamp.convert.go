package utils

import "time"

// ConvertTimestampToDate converts a Unix timestamp to a formatted date-time string.
func ConvertTimestampToDate(timestamp int64) string {
	// Convert timestamp to time.Time object
	t := time.Unix(timestamp, 0)

	// Format to a readable date-time format
	return t.Format("2006-01-02 15:04:05")
}
