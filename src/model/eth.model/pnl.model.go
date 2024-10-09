package ethmodel

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
	TotalETHBuy          float64 `json:"total-eth-buy" bson:"totalethbuy"`
	TotalETHSell         float64 `json:"total-eth-sell" bson:"totalethsell"`
	TotalETHSellActual   float64 `json:"total-eth-sell-actual" bson:"totalethsellactual"`
	TokenHoldAmount      float64 `json:"token-hold-amount" bson:"tokenholdamount"`
	TokenHoldETHAmount   float64 `json:"token-hold-eth-amount" bson:"tokenholdethamount"`
	ProfitETH            float64 `json:"profit-eth" bson:"profiteth"`
	ProfitETHActual      float64 `json:"profit-eth-actual" bson:"profitethactual"`
	XPNL                 float64 `json:"xpnl" bson:"xpnl"`
	XPNLRate             float64 `json:"xpnl-rate" bson:"xpnlrate"`
	PriceETHFirstBuy     float64 `json:"price-eth-first-buy" bson:"priceethfirstbuy"`
	PriceETHBestSell     float64 `json:"price-eth-best-sell" bson:"priceethbestsell"`
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
	TotalETHBuy          float64 `json:"total-eth-buy" bson:"totalethbuy"`
	TotalETHSell         float64 `json:"total-eth-sell" bson:"totalethsell"`
	TotalETHSellActual   float64 `json:"total-eth-sell-actual" bson:"totalethsellactual"`
	TokenHoldAmount      float64 `json:"token-hold-amount" bson:"tokenholdamount"`
	TokenHoldETHAmount   float64 `json:"token-hold-eth-amount" bson:"tokenholdethamount"`
	ProfitETH            float64 `json:"profit-eth" bson:"profiteth"`
	ProfitETHActual      float64 `json:"profit-eth-actual" bson:"profitethactual"`
	LostXPNL             float64 `json:"xpnl" bson:"lostxpnl"`
	LostXPNLRate         float64 `json:"lost-xpnl-rate" bson:"lostxpnlrate"`
	PriceETHFirstBuy     float64 `json:"price-eth-first-buy" bson:"priceethfirstbuy"`
	PriceETHBestSell     float64 `json:"price-eth-best-sell" bson:"priceethbestsell"`
	LostXPNLTrade        float64 `json:"xpnl-trade" bson:"xpnltrade"`
	LostXPNLRateTrade    float64 `json:"lost-xpnl-rate-trade" bson:"lostxpnlratetrade"`
	StartTime            string  `json:"start-time" bson:"starttime"`
	EndTime              string  `json:"end-time" bson:"endtime"`
}

type SummaryReview struct {
	TotalETHPNLAmount       float64 `json:"total-eth-pnl-amount" bson:"totalethpnlamount"`
	TotalETHPNLAmountActual float64 `json:"total-eth-pnl-amount-actual" bson:"totalethpnlamountactual"`
	TotalWin                int     `json:"total-win" bson:"totalwin"`
	TotalLost               int     `json:"total-lost" bson:"totallost"`
	WinRate                 float64 `json:"win-rate" bson:"winrate"`
	BigXPNL                 int     `json:"big-xpnl" bson:"bigxpnl"`
	RateBigXPNL             float64 `json:"rate-big-xpnl" bson:"ratebigxpnl"`
}

type EventTrade struct {
	Type         string  `json:"type" bson:"type"`
	EventType    string  `json:"event-type" bson:"eventtype"`
	PriceETH     float64 `json:"price-eth" bson:"priceeth"`
	ETHAmount    float64 `json:"eth-amount" bson:"ethamount"`
	TokensAmount float64 `json:"tokens-amount" bson:"tokensamount"`
	Timestamp    int64   `json:"timestamp" bson:"timestamp"`
	DateTime     string  `json:"date-time" bson:"datetime"`
}
