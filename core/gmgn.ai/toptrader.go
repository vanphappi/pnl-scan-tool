package gmgnai

import (
	"encoding/json"
	"fmt"
	"log"
)

type TopTradersData struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data []WalletData `json:"data"`
}

const baseUrlTrader = "https://gmgn.ai/defi/quotation/v1/tokens/top_traders"

func getTopTradersToken(chain string, token string) (*TopHoldersData, error) {
	var url string
	if chain == "sol" {
		url = fmt.Sprintf("%s/%s/%s?limit=%d&tag=All&orderby=realized_profit&direction=desc", baseUrlTrader, chain, token, limit)
	} else if chain == "eth" {
		url = fmt.Sprintf("%s/%s/%s?orderby=realized_profit&direction=desc", baseUrlTrader, chain, token)
	}

	fmt.Println(url)

	// if cursor != "" {
	// 	url += "&cursor=" + cursor
	// }

	// Fetch with retry logic
	result, err := fetchWithRetry(url)

	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var apiResponse TopHoldersData

	if err := json.Unmarshal(result, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func TopTradersToken(chain string, token string) []WalletData {
	var TopTraders []WalletData

	count := 0

	apiResponse, err := getTopTradersToken(chain, token)

	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	count += len(apiResponse.Data)

	fmt.Println("Scan Total Traders: ", count)

	// Append activities to the slice
	TopTraders = append(TopTraders, apiResponse.Data...)

	return TopTraders
}
