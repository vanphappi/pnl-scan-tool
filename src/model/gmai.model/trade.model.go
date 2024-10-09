package gmaimodel

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
