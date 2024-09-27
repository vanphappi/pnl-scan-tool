package geckoterminal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pnl-solana-tool/package/utils"
	"time"

	"golang.org/x/exp/rand"
)

// TokenHoldersData is the root structure containing the token holders details.
type TokenHoldersData struct {
	Data Data `json:"data"`
}

// Data holds the ID and the type of the token holder detail, along with attributes.
type Data struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

// Attributes contains the top holders' details.
type Attributes struct {
	TopHolders []Holder `json:"top_holders"`
}

// Holder represents the details of each top holder.
type Holder struct {
	WalletAddress   string  `json:"wallet_address"`
	Amount          float64 `json:"amount"`
	OwnerPercent    float64 `json:"owner_percent"`
	ValueInUSD      float64 `json:"value_in_usd"`
	IsLiquidityPool *bool   `json:"is_liquidity_pool"` // Use pointer to allow null values
}

func TopHolders(chain string, token string) (TokenHoldersData, error) {

	apiUrl := fmt.Sprintf("https://app.geckoterminal.com/api/p1/%s/tokens/%s/top_holders", chain, token)

	data, err := fetchDataGeckoTerminalFromAPI(apiUrl)

	if err != nil {
		return TokenHoldersData{}, err
	}

	var tokenHoldersData TokenHoldersData

	err = json.Unmarshal(data, &tokenHoldersData)

	if err != nil {
		return TokenHoldersData{}, err
	}

	return tokenHoldersData, nil
}

func fetchDataGeckoTerminalFromAPI(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Create a custom transport with TLS settings
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

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
