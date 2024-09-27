package photonsol

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"pnl-solana-tool/package/utils"
	"time"
)

// fetchDataPhotonFromAPI fetches data from the given API URL using randomized headers, cookies, and retry logic.
func fetchDataPhotonFromAPI(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Create a custom transport with TLS settings
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// proxyURL := proxy.GetRandomProxy()
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// 	Proxy:           http.ProxyURL(proxyURL),
	// }
	// client := &http.Client{Transport: tr}

	// Smart retry with exponential backoff and random jitter

	var body []byte

	maxRetries := 100

	for retries := 1; retries <= maxRetries; retries++ {
		// Set a random User-Agent
		req.Header.Set("User-Agent", utils.GetRandomUserAgent())

		req.Header.Set("Referer", url)

		// Set random headers
		headers := utils.GetRandomHeaders()
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		// Set random cookies
		cookies := utils.GetRandomCookies()

		for key, value := range cookies {
			req.AddCookie(&http.Cookie{Name: key, Value: value})
		}

		resp, err := client.Do(req)

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
