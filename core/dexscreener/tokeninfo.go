package dexscreener

import (
	"encoding/json"
	"fmt"
)

type TokenInfo struct {
	SchemaVersion string `json:"schemaVersion"`
	GP            struct {
		DataStatus           string `json:"dataStatus"`
		IsInDex              bool   `json:"isInDex"`
		IsOpenSource         bool   `json:"isOpenSource"`
		IsProxy              bool   `json:"isProxy"`
		IsTrueToken          bool   `json:"isTrueToken"`
		BuyTax               string `json:"buyTax"`
		CanTakeBackOwnership bool   `json:"canTakeBackOwnership"`
		CannotSellAll        bool   `json:"cannotSellAll"`
		CreatorAddress       string `json:"creatorAddress"`
		CreatorBalance       string `json:"creatorBalance"`
		CreatorPercent       string `json:"creatorPercent"`
		Dex                  []struct {
			Name      string `json:"name"`
			Liquidity string `json:"liquidity"`
			Pair      string `json:"pair"`
		} `json:"dex"`
		ExternalCall bool `json:"externalCall"`
		HolderCount  int  `json:"holderCount"`
		Holders      []struct {
			Address    string `json:"address"`
			IsLocked   bool   `json:"isLocked"`
			Tag        string `json:"tag,omitempty"`
			IsContract bool   `json:"isContract"`
			Balance    string `json:"balance"`
			Percent    string `json:"percent"`
		} `json:"holders"`
		IsAntiWhale         bool   `json:"isAntiWhale"`
		AntiWhaleModifiable bool   `json:"antiWhaleModifiable"`
		IsBlacklisted       bool   `json:"isBlacklisted"`
		IsHoneypot          bool   `json:"isHoneypot"`
		IsMintable          bool   `json:"isMintable"`
		IsWhitelisted       bool   `json:"isWhitelisted"`
		LpHolderCount       int    `json:"lpHolderCount"`
		LpTotalSupply       string `json:"lpTotalSupply"`
		LpHolders           []struct {
			Address    string `json:"address"`
			IsLocked   bool   `json:"isLocked"`
			IsContract bool   `json:"isContract"`
			Balance    string `json:"balance"`
			Percent    string `json:"percent"`
		} `json:"lpHolders"`
		OwnerAddress       string `json:"ownerAddress"`
		OwnerBalance       string `json:"ownerBalance"`
		OwnerChangeBalance bool   `json:"ownerChangeBalance"`
		OwnerPercent       string `json:"ownerPercent"`
		HiddenOwner        bool   `json:"hiddenOwner"`
		SellTax            string `json:"sellTax"`
		SlippageModifiable bool   `json:"slippageModifiable"`
		TotalSupply        string `json:"totalSupply"`
		TransferPausable   bool   `json:"transferPausable"`
		TradingCooldown    bool   `json:"tradingCooldown"`
		TokenName          string `json:"tokenName"`
		TokenSymbol        string `json:"tokenSymbol"`
		TrustList          bool   `json:"trustList"`
		UpdatedAt          int64  `json:"updatedAt"`
	} `json:"gp"`
	TS struct {
		Status       string   `json:"status"`
		ChainId      string   `json:"chainId"`
		Address      string   `json:"address"`
		TotalSupply  int      `json:"totalSupply"`
		DeployerAddr string   `json:"deployerAddr"`
		IsFlagged    bool     `json:"isFlagged"`
		Exploits     []string `json:"exploits"`
		Contract     struct {
			IsSourceVerified        bool `json:"isSourceVerified"`
			HasMint                 bool `json:"hasMint"`
			HasFeeModifier          bool `json:"hasFeeModifier"`
			HasMaxTransactionAmount bool `json:"hasMaxTransactionAmount"`
			HasBlocklist            bool `json:"hasBlocklist"`
			HasProxy                bool `json:"hasProxy"`
			HasPausable             bool `json:"hasPausable"`
		} `json:"contract"`
		Score       int    `json:"score"`
		RiskLevel   string `json:"riskLevel"`
		Permissions struct {
			OwnerAddress         string `json:"ownerAddress"`
			IsOwnershipRenounced bool   `json:"isOwnershipRenounced"`
		} `json:"permissions"`
		SwapSimulation struct {
			IsSellable bool `json:"isSellable"`
		} `json:"swapSimulation"`
		Balances struct {
			BurnBalance     int `json:"burnBalance"`
			LockBalance     int `json:"lockBalance"`
			DeployerBalance int `json:"deployerBalance"`
			OwnerBalance    int `json:"ownerBalance"`
			TopHolders      []struct {
				Address    string `json:"address"`
				Balance    int    `json:"balance"`
				IsContract bool   `json:"isContract"`
			} `json:"topHolders"`
		} `json:"balances"`
		Pools []interface{} `json:"pools"`
		Tests []struct {
			Id          string `json:"id"`
			Description string `json:"description"`
			Result      bool   `json:"result"`
		} `json:"tests"`
		UpdatedAt int64 `json:"updatedAt"`
	} `json:"ts"`
	QI struct {
		ContractVerified bool `json:"contractVerified"`
		TokenDetails     struct {
			TokenName        string `json:"tokenName"`
			TokenSymbol      string `json:"tokenSymbol"`
			TokenDecimals    int    `json:"tokenDecimals"`
			TokenOwner       string `json:"tokenOwner"`
			TokenSupply      int    `json:"tokenSupply"`
			TokenCreatedDate int64  `json:"tokenCreatedDate"`
			QuickiTokenHash  struct {
				ExactQHash   string `json:"exactQHash"`
				SimilarQHash string `json:"similarQHash"`
			} `json:"quickiTokenHash"`
		} `json:"tokenDetails"`
		TokenDynamicDetails struct {
			LastUpdatedTimestamp int64   `json:"lastUpdatedTimestamp"`
			IsHoneypot           bool    `json:"isHoneypot"`
			BuyTax               string  `json:"buyTax"`
			SellTax              string  `json:"sellTax"`
			TransferTax          string  `json:"transferTax"`
			MaxWallet            string  `json:"maxWallet"`
			MaxWalletPercent     string  `json:"maxWalletPercent"`
			LpPair               string  `json:"lpPair"`
			LpSupply             float64 `json:"lpSupply"`
			LpBurnedPercent      string  `json:"lpBurnedPercent"`
			PriceImpact          string  `json:"priceImpact"`
			Problem              bool    `json:"problem"`
		} `json:"tokenDynamicDetails"`
		QuickiAudit struct {
			ContractCreator            string   `json:"contractCreator"`
			ContractOwner              string   `json:"contractOwner"`
			ContractName               string   `json:"contractName"`
			ContractChain              string   `json:"contractChain"`
			ContractAddress            string   `json:"contractAddress"`
			ContractRenounced          bool     `json:"contractRenounced"`
			HiddenOwner                bool     `json:"hiddenOwner"`
			IsProxy                    bool     `json:"isProxy"`
			HasExternalContractRisk    bool     `json:"hasExternalContractRisk"`
			HasObfuscatedAddressRisk   bool     `json:"hasObfuscatedAddressRisk"`
			CanMint                    bool     `json:"canMint"`
			CanBurn                    bool     `json:"canBurn"`
			CanBlacklist               bool     `json:"canBlacklist"`
			CanMultiBlacklist          bool     `json:"canMultiBlacklist"`
			CanWhitelist               bool     `json:"canWhitelist"`
			CantWhitelistRenounced     bool     `json:"cantWhitelistRenounced"`
			CanUpdateFees              bool     `json:"canUpdateFees"`
			CantUpdateFeesRenounced    bool     `json:"cantUpdateFeesRenounced"`
			CanUpdateMaxWallet         bool     `json:"canUpdateMaxWallet"`
			CanUpdateMaxTx             bool     `json:"canUpdateMaxTx"`
			CanPauseTrading            bool     `json:"canPauseTrading"`
			CantPauseTradingRenounced  bool     `json:"cantPauseTradingRenounced"`
			HasTradingCooldown         bool     `json:"hasTradingCooldown"`
			CanUpdateWallets           bool     `json:"canUpdateWallets"`
			HasSuspiciousFunctions     bool     `json:"hasSuspiciousFunctions"`
			HasExternalFunctions       bool     `json:"hasExternalFunctions"`
			HasFeeWarning              bool     `json:"hasFeeWarning"`
			HasModifiedTransferWarning bool     `json:"hasModifiedTransferWarning"`
			ModifiedTransferFunctions  []string `json:"modifiedTransferFunctions"`
			SuspiciousFunctions        []string `json:"suspiciousFunctions"`
			HasScams                   bool     `json:"hasScams"`
			ContractLinks              []string `json:"contractLinks"`
			Functions                  []string `json:"functions"`
			OnlyOwnerFunctions         []string `json:"onlyOwnerFunctions"`
			MultiBlacklistFunctions    []string `json:"multiBlacklistFunctions"`
			HasGeneralVulnerabilities  bool     `json:"hasGeneralVulnerabilities"`
		} `json:"quickiAudit"`
		ChainId      string `json:"chainId"`
		TokenAddress string `json:"tokenAddress"`
		UpdatedAt    int64  `json:"updatedAt"`
	} `json:"qi"`
	LL struct {
		Locks []struct {
			Tag        string  `json:"tag"`
			Address    string  `json:"address"`
			Amount     string  `json:"amount"`
			Percentage float64 `json:"percentage"`
		} `json:"locks"`
		TotalPercentage float64 `json:"totalPercentage"`
	} `json:"ll"`
}

func TokenInfomation(chain string, tokenAddress string) (TokenInfo, error) {
	var tokenInfo TokenInfo

	apiUrl := fmt.Sprintf("https://io.dexscreener.com/dex/pair-details/v3/%s/%s", chain, tokenAddress)

	data, err := fetchWithRetry(apiUrl)

	if err != nil {
		return tokenInfo, err
	}

	err = json.Unmarshal(data, &tokenInfo)

	if err != nil {
		return tokenInfo, err
	}

	return tokenInfo, nil
}
