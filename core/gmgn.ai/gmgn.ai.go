package gmgnai

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"pnl-scan-tool/package/utils"
	"time"

	"golang.org/x/exp/rand"
)

const baseUrl = "https://gmgn.ai/defi/quotation/v1/wallet_activity/"
const limit = 100 // Maximize limit per request

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

	// Read headers from a JSON file
	config, err := utils.ReadHeadersFromFile("cookies/header/gmai.headers.json")
	if err != nil {
		fmt.Println("Error reading headers:", err)
		return nil, err
	}

	var body []byte
	// var err error
	for attempt := 0; attempt < maxRetries; attempt++ {

		// Set headers to mimic a real browser request
		req.Header.Set("accept", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Referer", "https://gmgn.ai/") // Set referer to appear like a browser visit
		req.Header.Set("Connection", "keep-alive")

		// // Set cookies
		// cookie := &http.Cookie{
		// 	Name:  "__cf_bm",
		// 	Value: "du18h5EJkRmsW3pbyydYrlKcCS5KnyQL_v1eBrq_aiU-1728406295-1.0.1.1-Hu4la6GZPD4iYPr8xFM2G_2MWhfxqCH4H2Oc.7hISuOOAnDtKF9qE7q6gGwLCdGgWRZP_B5IchEdvpaT2CDLug",
		// 	Path:  "/",
		// }

		// req.AddCookie(cookie)

		// Add cookies to the request
		for name, value := range config.Cookies {
			req.AddCookie(&http.Cookie{
				Name:  name,
				Value: value,
			})
		}

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

func ByPass(urls string) ([]byte, error) {
	urls = "http://203.29.19.253:8000/data?url=" + url.QueryEscape(urls)

	client := &http.Client{}

	method := "GET"

	req, err := http.NewRequest(method, urls, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil

}
