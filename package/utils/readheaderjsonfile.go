package utils

import (
	"encoding/json"
	"os"
)

// Define a struct to hold the JSON structure
type HeadersConfig struct {
	Cookies   map[string]string `json:"cookies"`
	UserAgent string            `json:"user_agent"`
}

// Function to read JSON headers from file
func ReadHeadersFromFile(filename string) (*HeadersConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config HeadersConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
