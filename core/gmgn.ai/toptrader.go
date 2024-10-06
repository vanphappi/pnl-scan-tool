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

const baseUrlTrader = "https://gmgn.ai/defi/quotation/v1/tokens/top_traders/sol/"

func getTopTradersToken(token string) (*TopHoldersData, error) {
	url := fmt.Sprintf("%s%s?limit=%d&tag=All&orderby=amount_percentage&direction=desc", baseUrlTrader, token, limit)
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

func TopTradersToken(token string) []WalletData {
	var TopTraders []WalletData

	count := 0

	apiResponse, err := getTopTradersToken(token)

	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	count += len(apiResponse.Data)

	fmt.Println("Scan Total Holder: ", count)

	// Append activities to the slice
	TopTraders = append(TopTraders, apiResponse.Data...)

	return TopTraders
}
