package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RenameCSVFilesToWallet renames all .csv files in the given directory to wallet.csv
func RenameCSVFilesToWallet(dir string) error {
	// Read all files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	// Iterate over all files
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Check if the file has a .csv extension
		if strings.HasSuffix(file.Name(), ".csv") {
			// Target path for renaming
			targetPath := filepath.Join(dir, "wallet.csv")

			// Check if wallet.csv already exists
			if _, err := os.Stat(targetPath); err == nil {
				fmt.Println("wallet.csv already exists. Skipping renaming.")
				continue
			}

			// Get the file's full path
			filePath := filepath.Join(dir, file.Name())

			// Rename the file
			err := os.Rename(filePath, targetPath)
			if err != nil {
				return fmt.Errorf("error renaming file %s: %w", file.Name(), err)
			}

			fmt.Printf("Renamed %s to wallet.csv\n", file.Name())
		}
	}

	return nil
}
