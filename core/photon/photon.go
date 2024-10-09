package photon

import (
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"pnl-scan-tool/package/utils"
	"time"
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
