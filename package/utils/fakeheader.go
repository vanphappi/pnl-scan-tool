package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/rand"
)

// List of user agents to mimic different browsers
var userAgents = []string{
	// Google Chrome on Windows
	"Mozilla/5.0 (Windows NT 11.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.126 Safari/537.36",

	// Safari on macOS
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_5_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15",

	// Firefox on Linux
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0",

	// Microsoft Edge on Windows
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:93.0) Gecko/20100101 Firefox/93.0",

	// Opera on macOS
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36 OPR/76.0.3965.181",

	// Google Chrome on Android
	"Mozilla/5.0 (Linux; Android 11; SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.5195.136 Mobile Safari/537.36",

	// Safari on iOS
	"Mozilla/5.0 (iPhone; CPU iPhone OS 15_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Mobile/15E148 Safari/604.1",

	// Internet Explorer 11
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; AS; rv:11.0) like Gecko",

	// Edge on macOS
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Version/15.0 Edge/99.0.1150.55",
}

// Random headers sets with additional variability
var headerSets = []map[string]string{
	{
		"Accept":                      "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":             "en-US,en;q=0.5",
		"Connection":                  "keep-alive",
		"Cache-Control":               "no-cache",
		"Priority":                    "u=1, i",
		"sec-ch-ua":                   "\"Not_A Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\"",
		"sec-ch-ua-arch":              "x86",
		"sec-ch-ua-bitness":           "64",
		"sec-ch-ua-full-version":      "101.0.4951.67",
		"sec-ch-ua-full-version-list": "\"Chromium\";v=\"101.0.4951.67\", \"Google Chrome\";v=\"101.0.4951.67\"",
		"sec-ch-ua-mobile":            "?0",
		"sec-ch-ua-model":             "",
		"sec-ch-ua-platform":          "Windows",
		"sec-ch-ua-platform-version":  "10.0",
		"sec-fetch-dest":              "document",
		"sec-fetch-mode":              "navigate",
		"sec-fetch-site":              "same-origin",
	},
	{
		"Accept":                      "application/json, text/plain, */*",
		"Accept-Language":             "en-GB,en;q=0.9",
		"Connection":                  "keep-alive",
		"Cache-Control":               "no-store",
		"Priority":                    "u=1, i",
		"sec-ch-ua":                   "\"Not_A Brand\";v=\"99\", \"Chromium\";v=\"128\", \"Google Chrome\";v=\"128\"",
		"sec-ch-ua-arch":              "arm",
		"sec-ch-ua-bitness":           "64",
		"sec-ch-ua-full-version":      "128.0.6613.120",
		"sec-ch-ua-full-version-list": "\"Chromium\";v=\"128.0.6613.120\", \"Google Chrome\";v=\"128.0.6613.120\"",
		"sec-ch-ua-mobile":            "?0",
		"sec-ch-ua-model":             "",
		"sec-ch-ua-platform":          "macOS",
		"sec-ch-ua-platform-version":  "14.6.1",
		"sec-fetch-dest":              "empty",
		"sec-fetch-mode":              "cors",
		"sec-fetch-site":              "same-origin",
	},
	{
		"Accept":                      "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":             "fr-FR,fr;q=0.9",
		"Connection":                  "close",
		"Cache-Control":               "max-age=0",
		"Priority":                    "u=2, i",
		"sec-ch-ua":                   "\"Not_A Brand\";v=\"99\", \"Chromium\";v=\"110\", \"Google Chrome\";v=\"110\"",
		"sec-ch-ua-arch":              "x64",
		"sec-ch-ua-bitness":           "64",
		"sec-ch-ua-full-version":      "110.0.5481.100",
		"sec-ch-ua-full-version-list": "\"Chromium\";v=\"110.0.5481.100\", \"Google Chrome\";v=\"110.0.5481.100\"",
		"sec-ch-ua-mobile":            "?1",
		"sec-ch-ua-model":             "Pixel 5",
		"sec-ch-ua-platform":          "Android",
		"sec-ch-ua-platform-version":  "12",
		"sec-fetch-dest":              "document",
		"sec-fetch-mode":              "navigate",
		"sec-fetch-site":              "cross-site",
	},
	{
		"Accept":                      "application/xml, text/xml, */*",
		"Accept-Language":             "es-ES,es;q=0.8",
		"Connection":                  "keep-alive",
		"Cache-Control":               "no-cache",
		"Priority":                    "u=1, i",
		"sec-ch-ua":                   "\"Not_A Brand\";v=\"99\", \"Chromium\";v=\"105\", \"Google Chrome\";v=\"105\"",
		"sec-ch-ua-arch":              "arm64",
		"sec-ch-ua-bitness":           "64",
		"sec-ch-ua-full-version":      "105.0.5195.102",
		"sec-ch-ua-full-version-list": "\"Chromium\";v=\"105.0.5195.102\", \"Google Chrome\";v=\"105.0.5195.102\"",
		"sec-ch-ua-mobile":            "?0",
		"sec-ch-ua-model":             "",
		"sec-ch-ua-platform":          "Linux",
		"sec-ch-ua-platform-version":  "5.15",
		"sec-fetch-dest":              "empty",
		"sec-fetch-mode":              "cors",
		"sec-fetch-site":              "same-origin",
	},
	{
		"Accept":                      "application/json, text/html, */*",
		"Accept-Language":             "de-DE,de;q=0.7",
		"Connection":                  "keep-alive",
		"Cache-Control":               "no-cache",
		"Priority":                    "u=1, i",
		"sec-ch-ua":                   "\"Not_A Brand\";v=\"99\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"",
		"sec-ch-ua-arch":              "x86",
		"sec-ch-ua-bitness":           "32",
		"sec-ch-ua-full-version":      "120.0.6000.180",
		"sec-ch-ua-full-version-list": "\"Chromium\";v=\"120.0.6000.180\", \"Google Chrome\";v=\"120.0.6000.180\"",
		"sec-ch-ua-mobile":            "?0",
		"sec-ch-ua-model":             "",
		"sec-ch-ua-platform":          "Windows",
		"sec-ch-ua-platform-version":  "11.0",
		"sec-fetch-dest":              "document",
		"sec-fetch-mode":              "navigate",
		"sec-fetch-site":              "same-origin",
	},
}

// Cookies list to rotate during requests
var cookieSets = ReadCookiesFromFile()

// getRandomUserAgent returns a random User-Agent string
func GetRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// getRandomHeaders returns a random set of headers from headerSets
func GetRandomHeaders() map[string]string {
	return headerSets[rand.Intn(len(headerSets))]
}

// getRandomCookies returns a random set of cookies from cookieSets
func GetRandomCookies() map[string]string {
	return cookieSets[rand.Intn(len(cookieSets))]
}

func ReadCookiesFromFile() []map[string]string {
	var cookies []map[string]string
	count := 1
	for {
		// File path to your cookies file
		filePath := fmt.Sprintf("cookies/cookie%d.txt", count)

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			// fmt.Println("Error opening file:", err)
			return cookies
		}
		defer file.Close()

		// Read the entire file content
		scanner := bufio.NewScanner(file)
		var fileContent string
		for scanner.Scan() {
			fileContent += scanner.Text() + "\n"
		}
		if err := scanner.Err(); err != nil {
			// fmt.Println("Error reading file:", err)
			return cookies
		}

		// Define the cookies map
		cookie := map[string]string{
			"ws_sol_production":      "",
			"trid":                   "",
			"_photon_eth_production": "",
			"__cf_bm":                "",
			"cf_clearance":           "",
			"_photon_ta":             "",
		}

		// Split the file content by semicolons
		parts := strings.Split(fileContent, ";")

		// Iterate over the parts to populate the cookies map
		for _, part := range parts {
			// Trim leading and trailing spaces
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			// Split each part by '=' to get key and value
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			// Update the map if the key is present
			if _, exists := cookie[key]; exists {
				cookie[key] = value
			}
		}

		cookies = append(cookies, cookie)
		count++
	}
}
