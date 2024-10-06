package dexscreener

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

var client = &http.Client{
	Timeout: 30 * time.Second, // Set a timeout for the entire request
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true, // Prioritize server's cipher suite order
			CurvePreferences: []tls.CurveID{
				tls.CurveP256, tls.X25519, // Strong elliptic curves
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
		},
		MaxIdleConns:        100,              // Pool idle connections
		MaxIdleConnsPerHost: 10,               // Per-host connection limit
		IdleConnTimeout:     90 * time.Second, // Timeout for idle connections
		MaxConnsPerHost:     20,               // Limit the maximum simultaneous connections to a host
		DisableKeepAlives:   false,            // Enable keep-alive for better performance
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // Timeout for dialing the connection
			KeepAlive: 30 * time.Second, // Keep-alive settings
		}).DialContext,
		ExpectContinueTimeout: 1 * time.Second,           // Optimize for HTTP/1.1 Continue handling
		TLSHandshakeTimeout:   10 * time.Second,          // Set timeout for TLS handshake
		ResponseHeaderTimeout: 10 * time.Second,          // Set a timeout for waiting on response headers
		Proxy:                 http.ProxyFromEnvironment, // Use system-wide proxy settings
	},
}

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
