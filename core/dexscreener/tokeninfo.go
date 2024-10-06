package dexscreener

import (
	"encoding/json"
	"fmt"
)

type Chain struct {
	ID string `json:"id"`
}

type Website struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Social struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Profile struct {
	Header    bool   `json:"header"`
	Website   bool   `json:"website"`
	Twitter   bool   `json:"twitter"`
	Discord   bool   `json:"discord"`
	LinkCount int    `json:"linkCount"`
	ImgKey    string `json:"imgKey"`
}

type Token struct {
	ID              string    `json:"id"`
	Chain           Chain     `json:"chain"`
	Address         string    `json:"address"`
	Name            string    `json:"name"`
	Symbol          string    `json:"symbol"`
	Description     string    `json:"description"`
	Websites        []Website `json:"websites"`
	Socials         []Social  `json:"socials"`
	LockedAddresses []string  `json:"lockedAddresses"`
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string    `json:"updatedAt"`
	SortByDate      string    `json:"sortByDate"`
	Image           string    `json:"image"`
	HeaderImage     string    `json:"headerImage"`
	Profile         Profile   `json:"profile"`
}

type TokenInfo struct {
	SchemaVersion string      `json:"schemaVersion"`
	CG            *string     `json:"cg"`
	GP            *string     `json:"gp"`
	TS            *string     `json:"ts"`
	CMC           *string     `json:"cmc"`
	TI            Token       `json:"ti"`
	CMS           *string     `json:"cms"`
	QI            interface{} `json:"qi"`
	DS            Token       `json:"ds"`
	LL            *string     `json:"ll"`
	Holders       *string     `json:"holders"`
	LPHolders     *string     `json:"lpHolders"`
	SU            *string     `json:"su"`
	TA            *string     `json:"ta"`
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
