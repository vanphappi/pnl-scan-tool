package solmodel

type PNL struct {
	WalletAddress string         `json:"wallet-address" bson:"walletaddress"`
	TradeHistory  []TradeHistory `json:"trades" bson:"trades"`
	XPNLs         []XPNL         `json:"xpnl" bson:"xpnl"`
	LostXPNLs     []LostXPNL     `json:"lost-xpnl" bson:"lostxpnl"`
	SummaryReview SummaryReview  `json:"summary-review" bson:"summaryreview"`
}

type TradeHistory struct {
	TokenAddress string       `json:"token-address" bson:"tokenaddress"`
	TokenSymbol  string       `json:"token-symbol" bson:"tokensymbol"`
	EventTrades  []EventTrade `json:"event-trades" bson:"eventtrades"`
	StartTime    string       `json:"start-time" bson:"starttime"`
	EndTime      string       `json:"end-time" bson:"endtime"`
}

type XPNL struct {
	TokenAddress         string  `json:"token-address" bson:"tokenaddress"`
	TokenSymbol          string  `json:"token-symbol" bson:"tokensymbol"`
	CountBuy             int     `json:"count-buy" bson:"countbuy"`
	CountSell            int     `json:"count-sell" bson:"countsell"`
	CountSellActual      int     `json:"count-sell-actual" bson:"countsellactual"`
	TotalTokenBuy        float64 `json:"total-token-buy" bson:"totaltokenbuy"`
	TotalTokenSell       float64 `json:"total-token-sell" bson:"totaltokensell"`
	TotalTokenSellActual float64 `json:"total-token-sell-actual" bson:"totaltokensellactual"`
	TotalSolBuy          float64 `json:"total-sol-buy" bson:"totalsolbuy"`
	TotalSolSell         float64 `json:"total-sol-sell" bson:"totalsolsell"`
	TotalSolSellActual   float64 `json:"total-sol-sell-actual" bson:"totalsolsellactual"`
	TokenHoldAmount      float64 `json:"token-hold-amount" bson:"tokenholdamount"`
	TokenHoldSolAmount   float64 `json:"token-hold-sol-amount" bson:"tokenholdsolamount"`
	ProfitSol            float64 `json:"profit-sol" bson:"profitsol"`
	ProfitSolActual      float64 `json:"profit-sol-actual" bson:"profitsolactual"`
	XPNL                 float64 `json:"xpnl" bson:"xpnl"`
	XPNLRate             float64 `json:"xpnl-rate" bson:"xpnlrate"`
	PriceSolFirstBuy     float64 `json:"price-sol-first-buy" bson:"pricesolfirstbuy"`
	PriceSolBestSell     float64 `json:"price-sol-best-sell" bson:"pricesolbestsell"`
	XPNLTrade            float64 `json:"xpnl-trade" bson:"xpnltrade"`
	XPNLRateTrade        float64 `json:"xpnl-rate-trade" bson:"xpnlratetrade"`
	StartTime            string  `json:"start-time" bson:"starttime"`
	EndTime              string  `json:"end-time" bson:"endtime"`
}

type LostXPNL struct {
	TokenAddress         string  `json:"token-address" bson:"tokenaddress"`
	TokenSymbol          string  `json:"token-symbol" bson:"tokensymbol"`
	CountBuy             int     `json:"count-buy" bson:"countbuy"`
	CountSell            int     `json:"count-sell" bson:"countsell"`
	CountSellActual      int     `json:"count-sell-actual" bson:"countsellactual"`
	TotalTokenBuy        float64 `json:"total-token-buy" bson:"totaltokenbuy"`
	TotalTokenSell       float64 `json:"total-token-sell" bson:"totaltokensell"`
	TotalTokenSellActual float64 `json:"total-token-sell-actual" bson:"totaltokensellactual"`
	TotalSolBuy          float64 `json:"total-sol-buy" bson:"totalsolbuy"`
	TotalSolSell         float64 `json:"total-sol-sell" bson:"totalsolsell"`
	TotalSolSellActual   float64 `json:"total-sol-sell-actual" bson:"totalsolsellactual"`
	TokenHoldAmount      float64 `json:"token-hold-amount" bson:"tokenholdamount"`
	TokenHoldSolAmount   float64 `json:"token-hold-sol-amount" bson:"tokenholdsolamount"`
	ProfitSol            float64 `json:"profit-sol" bson:"profitsol"`
	ProfitSolActual      float64 `json:"profit-sol-actual" bson:"profitsolactual"`
	LostXPNL             float64 `json:"xpnl" bson:"lostxpnl"`
	LostXPNLRate         float64 `json:"lost-xpnl-rate" bson:"lostxpnlrate"`
	PriceSolFirstBuy     float64 `json:"price-sol-first-buy" bson:"pricesolfirstbuy"`
	PriceSolBestSell     float64 `json:"price-sol-best-sell" bson:"pricesolbestsell"`
	LostXPNLTrade        float64 `json:"xpnl-trade" bson:"xpnltrade"`
	LostXPNLRateTrade    float64 `json:"lost-xpnl-rate-trade" bson:"lostxpnlratetrade"`
	StartTime            string  `json:"start-time" bson:"starttime"`
	EndTime              string  `json:"end-time" bson:"endtime"`
}

type SummaryReview struct {
	TotalSolPNLAmount       float64 `json:"total-sol-pnl-amount" bson:"totalsolpnlamount"`
	TotalSolPNLAmountActual float64 `json:"total-sol-pnl-amount-actual" bson:"totalsolpnlamountactual"`
	TotalWin                int     `json:"total-win" bson:"totalwin"`
	TotalLost               int     `json:"total-lost" bson:"totallost"`
	WinRate                 float64 `json:"win-rate" bson:"winrate"`
	BigXPNL                 int     `json:"big-xpnl" bson:"bigxpnl"`
	RateBigXPNL             float64 `json:"rate-big-xpnl" bson:"ratebigxpnl"`
}

type EventTrade struct {
	Type         string  `json:"type" bson:"type"`
	EventType    string  `json:"event-type" bson:"eventtype"`
	PriceSol     float64 `json:"price-sol" bson:"pricesol"`
	SolAmount    float64 `json:"sol-amount" bson:"solamount"`
	TokensAmount float64 `json:"tokens-amount" bson:"tokensamount"`
	Timestamp    int64   `json:"timestamp" bson:"timestamp"`
	DateTime     string  `json:"date-time" bson:"datetime"`
}
