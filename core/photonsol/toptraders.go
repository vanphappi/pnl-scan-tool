package photonsol

import (
	"encoding/json"
	"fmt"
)

// Attributes represents the attributes of a transaction
type AttributesTopTrader struct {
	BoughtUsd   string `json:"boughtUsd"`
	BoughtToken string `json:"boughtToken"`
	BoughtCount int    `json:"boughtCount"`
	SoldUsd     string `json:"soldUsd"`
	// SoldToken       string      `json:"soldToken"`
	SoldCount       int         `json:"soldCount"`
	RemainingTokens interface{} `json:"remainingTokens"`
	PlUsd           string      `json:"plUsd"`
	Kind            string      `json:"kind"`
	Signer          string      `json:"signer"`
}

// Transaction represents the overall transaction structure
type TopTrader struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes AttributesTopTrader `json:"attributes"`
}

// Define a new struct to match the JSON response structure
type TopTraders struct {
	Data []TopTrader `json:"data"`
}

func (t *Token) TopTraders(poolId int) ([]TopTrader, error) {
	apiUrl := fmt.Sprintf("https://photon-sol.tinyastro.io/api/events/top_traders?order_by=timestamp&order_dir=dir&pool_id=%d&page=1", poolId)

	data, err := fetchDataPhotonFromAPI(apiUrl)

	if err != nil {
		return nil, err
	}

	var topTraders TopTraders

	err = json.Unmarshal(data, &topTraders)

	if err != nil {
		return nil, err
	}

	return topTraders.Data, nil
}
