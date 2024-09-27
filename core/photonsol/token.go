package photonsol

import (
	"encoding/json"
	"fmt"
	"pnl-solana-tool/package/utils"
	"strings"
)

type Token struct {
	TokenAddress string
}

type TokenInfomation struct {
	// AmountIndex         int      `json:"amount-index"`
	// PoolAddress         string   `json:"pool-address"`
	PriceUsd   float64 `json:"price-usd"`
	PriceQuote float64 `json:"price-quote"`
	// CreatedAt           int64    `json:"created_at"`
	PoolId  int `json:"pool-id"`
	TokenId int `json:"token-id"`
	// Dex                 string   `json:"dex"`
	// FromPump            bool     `json:"from-pump"`
	// PumpPoolId          *int     `json:"pump-pool_id"`
	// PoolFactoryName     string   `json:"pool-factory-name"`
	// Decimals            int      `json:"decimals"`
	TokenAddress string `json:"token-address"`
	TokenSymbol  string `json:"token-symbol"`
	// HumanTotalSupply    float64  `json:"human-total-supply"`
	// TotalSupply         *big.Int `json:"total-supply"`
	// DefaultInterval     string   `json:"default-interval"`
	// FirstEventTimestamp int64    `json:"first-event-timestamp"`
	// LastEventTimestamp  int64    `json:"last-event-timestamp"`
	// PurchasesCount      int      `json:"purchases-count"`
	// ChartHeight         string   `json:"chart-height"`
	// TxSort              string   `json:"tx-sort"`
	// TxDateAge           string   `json:"tx-date-age"`
	// IsWatching          bool     `json:"is_watching"`
	// OwnerWallet         string   `json:"owner-wallet"`
	// DeployerWallet      string   `json:"deployer-wallet"`
	// CaDeployerWallet    string   `json:"ca-deployer-wallet"`
	// LpDeployerWallet    string   `json:"lp-deployer-wallet"`
	// Outliers            bool     `json:"outliers"`
	// MyTrades            bool     `json:"my-trades"`
	// DevTrades           bool     `json:"dev-trades"`
	// OtherTrades         []string `json:"other-trades"`
	// AvgPriceLine        bool     `json:"avg-price-line"`
	// OrderLines          bool     `json:"order-lines"`
	// ChartCurrency       string   `json:"chart-currency"`
	// ChartType           string   `json:"chart-type"`
	// IsInternalWallet    bool     `json:"is-internal-wallet"`
	// QuoteToken0         bool     `json:"quote-token0"`
	// Reserve0            *big.Int `json:"reserve0"`
	// Reserve1            *big.Int `json:"reserve1"`
	// Reserve0Tr          *big.Int `json:"reserve0Tr"`
	// Reserve1Tr          *big.Int `json:"reserve1Tr"`
	// Minliq              int      `json:"minliq"`
	// JitoTips            float64  `json:"jito-tips"`
	// MigratedAddress     *string  `json:"migrated-address"`
	// Ignored             bool     `json:"ignored"`
	// OldPool             bool     `json:"old-pool"`
	// PoolAddresses       []string `json:"pool-addresses"`
	// DefaultBribery      float64  `json:"default-bribery"`
	// Snipers             *string  `json:"snipers"`
	// Insiders            *string  `json:"insiders"`
	// Tax                 float64  `json:"tax"`
	// Is2022              bool     `json:"is2022"`
}

// TokenInfomation fetches the token information from the Photon website.
//
// The returned struct contains the token information, which includes the token
// address, token symbol, price in USD, price in quote token, pool ID, token ID,
// and other relevant information.
//
// The function returns an error if it fails to fetch the data or parse the JSON
// response.
func (p *Token) TokenInfomation() (*TokenInfomation, error) {
	apiUrl := fmt.Sprintf("https://photon-sol.tinyastro.io/en/lp/%s", p.TokenAddress)

	result, err := fetchDataPhotonFromAPI(apiUrl)

	data := string(result)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}

	// Extract the relevant JSON data from the JavaScript
	jsonStart := strings.Index(data, "window.taConfig.show = ")

	if jsonStart == -1 {
		return nil, fmt.Errorf("couldn't find window.taConfig.show in the response")
	}
	jsonStart += len("window.taConfig.show = ")

	jsonEnd := strings.Index(data[jsonStart:], "};")

	if jsonEnd == -1 {
		return nil, fmt.Errorf("couldn't find the end of the JSON object")
	}
	jsonEnd += jsonStart + 1

	jsonStr := data[jsonStart:jsonEnd]

	// Preprocess the JSON string
	jsonStr = utils.PreprocessJSON(jsonStr)
	jsonStr = preprocessJSON(jsonStr)

	// Parse the JSON data
	var tokenInfo TokenInfomation

	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	decoder.UseNumber()

	err = json.Unmarshal([]byte(jsonStr), &tokenInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v\nJSON string: %s", err, jsonStr)
	}

	return &tokenInfo, nil
}

func preprocessJSON(jsonStr string) string {
	var result strings.Builder
	inTokenSymbol := false
	inQuotes := false
	tokenSymbolValue := ""

	for i := 0; i < len(jsonStr); i++ {
		char := jsonStr[i]

		if !inTokenSymbol {
			if strings.HasPrefix(jsonStr[i:], `"token-symbol"`) {
				inTokenSymbol = true
				result.WriteString(`"token-symbol":`)
				i += len(`"token-symbol"`)
				continue
			}
			result.WriteByte(char)
		} else {
			if char == ':' && !inQuotes {
				continue
			}
			if char == ',' || char == '}' {
				// End of token-symbol value
				escapedValue := strings.Replace(tokenSymbolValue, `"`, `\"`, -1)
				result.WriteString(fmt.Sprintf(`"%s"`, escapedValue))
				result.WriteByte(char)
				inTokenSymbol = false
				tokenSymbolValue = ""
				continue
			}
			if char == '"' && (i == 0 || jsonStr[i-1] != '\\') {
				inQuotes = !inQuotes
			}
			tokenSymbolValue += string(char)
		}
	}

	return result.String()
}
