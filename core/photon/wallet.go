package photon

import (
	"encoding/json"
	"fmt"
)

type Wallet struct {
	WalletAddress string
}

type Transaction struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Timestamp    int64       `json:"timestamp"`
	Type         string      `json:"type"`
	EventType    string      `json:"eventType"`
	PriceUsd     string      `json:"priceUsd"`
	PriceQuote   string      `json:"priceQuote"`
	UsdAmount    interface{} `json:"usdAmount"`
	QuoteAmount  string      `json:"quoteAmount"`
	TokensAmount string      `json:"tokensAmount"`
	Amt0         string      `json:"amt0"`
	Amt1         string      `json:"amt1"`
	Maker        string      `json:"maker"`
	TxHash       string      `json:"txHash"`
	IsOutlier    bool        `json:"isOutlier"`
	SortId       int64       `json:"sortId"`
	Kind         string      `json:"kind"`
	Slot         int64       `json:"slot"`
}

// Define a new struct to match the JSON response structure
type Transactions struct {
	Data []Transaction `json:"data"`
}

// Transactions fetches the transactions of the given wallet address from Photon using the Photon API.
//
// It takes an optional parameter `poolId` which is the ID of the pool to filter by.
// If `poolId` is 0, it fetches all transactions.
//
// It will then fetch the token information and the top traders of the token, and for each top trader, it will fetch the transactions and calculate the PLN history.
// It will then check if the PLN history meets certain conditions, and if so, it will append the wallet address to a file named "wallet.pnl.txt".
// It will then insert the token information into the database.
//
// The function returns a list of Transaction objects which represent the transactions of the given wallet address.
// If the function fails, it will return an error.
func (w *Wallet) Transactions(poolId int) ([]Transaction, error) {
	apiUrl := fmt.Sprintf("https://photon-sol.tinyastro.io/api/lp/events?old_pool=false&order_by=timestamp&order_dir=desc&pool_id=%d&signer=%s", poolId, w.WalletAddress)
	data, err := fetchDataPhotonFromAPI(apiUrl)

	if err != nil {
		return nil, err
	}

	var transactions Transactions

	err = json.Unmarshal(data, &transactions)

	if err != nil {
		return nil, err
	}

	return transactions.Data, nil
}
