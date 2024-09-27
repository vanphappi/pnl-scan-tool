package files

import (
	"fmt"
	"os"
)

// appendToFile appends the given content to the specified file
func AppendToFile(filename string, content string) error {
	// Open file in append mode, create it if it doesn't exist, write only
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file with a newline at the end
	_, err = file.WriteString(content + "\n")
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile deletes the file at the provided path
func DeleteFile(filePath string) error {
	// Use os.Remove to delete the file
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}
	fmt.Printf("File %s deleted successfully\n", filePath)
	return nil
}

// FileExists checks if a file exists and is not a directory
func FileExists(filePath string) bool {
	// Use os.Stat to get the file info
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// File does not exist
		return false
	}
	// Check if it's not a directory
	return !info.IsDir()
}
