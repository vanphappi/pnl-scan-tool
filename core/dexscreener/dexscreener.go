package dexscreener

import (
	"fmt"
	"io"
	"net/http"
	"pnl-scan-tool/package/utils"
	"time"

	"golang.org/x/exp/rand"
)

const maxRetries = 1000 // Maximum retry attempts

// Fetch data from the API with retries and exponential backoff
func fetchWithRetry(url string) ([]byte, error) {
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	var body []byte
	// var err error
	for attempt := 0; attempt < maxRetries; attempt++ {

		// Set headers to mimic a real browser request
		req.Header.Set("accept", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Referer", "https://io.dexscreener.com/") // Set referer to appear like a browser visit
		req.Header.Set("Connection", "keep-alive")

		// Set cookies
		// Add new cookies
		req.Header.Add("Cookie", "__cf_bm=y4URIekeGY_vAVluVYHE2NEFk_9zrh0WK7iWEogWbmM-1728135408-1.0.1.1-vdfk6WRLaecbO1dLAgMjvuuJoD2Zt7yMhaP8lV6NiwuK0igE5Q8JF1mjNeCYY2xDJozOcT3ageImOGPJadog1FzJmnHBliK5MS_umJm01MI")
		req.Header.Add("Cookie", "__cflb=0H28vzQ7jjUXq92cxrCqHJ17hceAH2AYkTHPrAUjHWV")

		// Add new User-Agent
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
		// Make the request
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
			// Exponential backoff with jitter
			baseWaitTime := time.Duration(2<<attempt) * time.Second
			jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
			totalWaitTime := baseWaitTime + jitter
			fmt.Printf("Received %d, retrying in %v...\n", resp.StatusCode, totalWaitTime)
			fmt.Println("Respone:", resp.Status)
			// time.Sleep(totalWaitTime)

			// respBody, err := ByPass(url)

			// // fmt.Println("Request:", url)

			// if err != nil {
			// 	continue
			// }

			// return respBody, err

			// time.Sleep(1 * time.Second)

			// proxyURL = proxy.GetRandomProxy()
			// tr.Proxy = http.ProxyURL(proxyURL)

			continue
		}

		// Handle other non-successful responses
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed with status %d, body: %s\n", resp.StatusCode, respBody)
	}

	return nil, err
}
