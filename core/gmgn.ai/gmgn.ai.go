package gmgnai

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

const baseUrl = "https://gmgn.ai/defi/quotation/v1/wallet_activity/sol"
const limit = 100 // Maximize limit per request

// Response structure (example, adjust based on actual API)
type ApiResponseGMGNAI struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

type Data struct {
	Activities []Activity `json:"activities"`
	Next       string     `json:"next"`
}

type Activity struct {
	Chain          string  `json:"chain"`
	TxHash         string  `json:"tx_hash"`
	Timestamp      int64   `json:"timestamp"`
	EventType      string  `json:"event_type"`
	TokenAddress   string  `json:"token_address"`
	Token          Token   `json:"token"`
	TokenAmount    string  `json:"token_amount"`
	QuoteAmount    string  `json:"quote_amount"`
	CostUSD        float64 `json:"cost_usd"`
	BuyCostUSD     float64 `json:"buy_cost_usd"`
	Price          float64 `json:"price"`
	PriceUSD       float64 `json:"price_usd"`
	IsOpenOrClose  int     `json:"is_open_or_close"`
	QuoteAddress   string  `json:"quote_address"`
	QuoteToken     Token   `json:"quote_token"`
	FromAddress    *string `json:"from_address"`     // Nullable field
	FromIsContract *bool   `json:"from_is_contract"` // Nullable field
	ToAddress      *string `json:"to_address"`       // Nullable field
	ToIsContract   *bool   `json:"to_is_contract"`   // Nullable field
	Balance        string  `json:"balance"`
	Sell30d        int     `json:"sell_30d"`
	LastActiveTs   int64   `json:"last_active_timestamp"`
	ID             string  `json:"id"`
}

type Token struct {
	Address  string  `json:"address"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Decimals int     `json:"decimals"`
	Logo     string  `json:"logo"`
	Price    float64 `json:"price"`
}

// Create a new HTTP client with a 30-second timeout
var client = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // Enable certificate verification
		},
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
		req.Header.Set("Referer", "https://solscan.io/") // Set referer to appear like a browser visit
		req.Header.Set("Connection", "keep-alive")

		// Set cookies
		cookie := &http.Cookie{
			Name:  "__cf_bm",
			Value: "THddocRJFL1aPybdb7rmv_rq4qq85i1zhmaotlBwmoI-1727366102-1.0.1.1-ibFVwSntsxA_gdnnFwcIPI9qcuiyJp4zBDSXOSrN_hbVTZpORpLqk8on8Mla3DtlgWhzVtRieJZ6M4G6mdaz4Q",
			Path:  "/",
		}

		req.AddCookie(cookie)

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
			// time.Sleep(totalWaitTime)

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
