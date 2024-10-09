package photon

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"pnl-scan-tool/package/utils"
	"time"
)

// fetchDataPhotonFromAPI fetches data from the given API URL using randomized headers, cookies, and retry logic.
func fetchDataPhotonFromAPI(url string) ([]byte, error) {
	// Read headers from a JSON file
	config, err := utils.ReadHeadersFromFile("cookies/header/photon.headers.json")
	if err != nil {
		fmt.Println("Error reading headers:", err)
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add cookies to the request
	for name, value := range config.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", config.UserAgent)

	var body []byte

	maxRetries := 100

	for retries := 1; retries <= maxRetries; retries++ {
		// // Set a random User-Agent
		// req.Header.Set("User-Agent", utils.GetRandomUserAgent())

		// req.Header.Set("Referer", url)

		// // Set random headers
		// headers := utils.GetRandomHeaders()
		// for key, value := range headers {
		// 	req.Header.Set(key, value)
		// }

		// // Set random cookies
		// cookies := utils.GetRandomCookies()

		// for key, value := range cookies {
		// 	req.AddCookie(&http.Cookie{Name: key, Value: value})
		// }

		resp, err := utils.Client.Do(req)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch data: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err = io.ReadAll(resp.Body)

			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}
			return body, nil
		}

		// Handle Too Many Requests (429) or Forbidden (403) with dynamic backoff
		if resp.StatusCode == 429 || resp.StatusCode == 403 {
			fmt.Println(url)
			// Exponential backoff with jitter
			baseWaitTime := time.Duration(2<<retries) * time.Second
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			totalWaitTime := baseWaitTime + jitter
			fmt.Printf("Received %d, retrying in %v...\n", resp.StatusCode, totalWaitTime)
			time.Sleep(totalWaitTime)
			// proxyURL = proxy.GetRandomProxy()
			// tr.Proxy = http.ProxyURL(proxyURL)
			continue
		}

		// Handle other non-successful responses
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed with status %d, body: %s\n", resp.StatusCode, respBody)
	}

	return nil, fmt.Errorf("exceeded maximum retries")
}
