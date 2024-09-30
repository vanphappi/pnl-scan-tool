package gmgnai

import (
	"encoding/json"
	"fmt"
	"log"
)

type TagRank struct {
	BluechipOwner int `json:"bluechip_owner,omitempty"`
}

type WalletData struct {
	Address           string   `json:"address"`
	AddrType          int      `json:"addr_type"`
	AmountCur         float64  `json:"amount_cur"`
	UsdValue          float64  `json:"usd_value"`
	CostCur           float64  `json:"cost_cur"`
	SellAmountCur     float64  `json:"sell_amount_cur"`
	SellAmountPercent float64  `json:"sell_amount_percentage"`
	SellVolumeCur     float64  `json:"sell_volume_cur"`
	BuyVolumeCur      float64  `json:"buy_volume_cur"`
	BuyAmountCur      float64  `json:"buy_amount_cur"`
	NetflowUsd        float64  `json:"netflow_usd"`
	NetflowAmount     float64  `json:"netflow_amount"`
	BuyTxCountCur     int      `json:"buy_tx_count_cur"`
	SellTxCountCur    int      `json:"sell_tx_count_cur"`
	WalletTagV2       string   `json:"wallet_tag_v2"`
	EthBalance        string   `json:"eth_balance"`
	SolBalance        string   `json:"sol_balance"`
	TrxBalance        string   `json:"trx_balance"`
	Balance           string   `json:"balance"`
	Profit            float64  `json:"profit"`
	RealizedProfit    float64  `json:"realized_profit"`
	UnrealizedProfit  float64  `json:"unrealized_profit"`
	ProfitChange      *float64 `json:"profit_change"`
	AmountPercentage  float64  `json:"amount_percentage"`
	AvgCost           *float64 `json:"avg_cost"`
	AvgSold           *float64 `json:"avg_sold"`
	Tags              []string `json:"tags"`
	MakerTokenTags    []string `json:"maker_token_tags"`
	Name              *string  `json:"name"`
	Avatar            *string  `json:"avatar"`
	TwitterUsername   *string  `json:"twitter_username"`
	TwitterName       *string  `json:"twitter_name"`
	TagRank           TagRank  `json:"tag_rank,omitempty"`
	LastActiveTS      int64    `json:"last_active_timestamp"`
	AccuAmount        float64  `json:"accu_amount"`
	AccuCost          float64  `json:"accu_cost"`
	Cost              float64  `json:"cost"`
	TotalCost         float64  `json:"total_cost"`
}

type TopHoldersData struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data []WalletData `json:"data"`
}

const baseUrlTopHolder = "https://gmgn.ai/defi/quotation/v1/tokens/top_holders/sol/"

func getTopHoldersToken(token string) (*TopHoldersData, error) {
	url := fmt.Sprintf("%s%s?limit=%d&tag=All&orderby=amount_percentage&direction=desc", baseUrlTopHolder, token, limit)
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

func TopHoldersToken(token string) []WalletData {
	var TopTraders []WalletData

	count := 0

	apiResponse, err := getTopHoldersToken(token)

	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	count += len(apiResponse.Data)

	fmt.Println("Scan Total Holder: ", count)

	// Append activities to the slice
	TopTraders = append(TopTraders, apiResponse.Data...)

	return TopTraders
}
