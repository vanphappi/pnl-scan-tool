package services

import (
	"fmt"
	"pnl-scan-tool/core/dexscreener"
	gmgnai "pnl-scan-tool/core/gmgn.ai"
	"pnl-scan-tool/core/photon"
	"pnl-scan-tool/package/files"
	"pnl-scan-tool/platform/database/mongodb"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func TopTraderScan(chain string, tokenAddress string) {

	if chain == "sol" {
		filter := bson.M{"tokenaddress": tokenAddress}

		_, err := mongodb.FindOne("token_scan_sol", filter)

		if err == nil {
			fmt.Println("Token address already exists in the database.")
			return
		}

		token := photon.Token{
			TokenAddress: tokenAddress,
		}

		data, err := token.TokenInfomation()

		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}

		// topTraders, err := token.TopTraders(data.PoolId)

		// if err != nil {
		// 	fmt.Println("Error: " + err.Error())
		// 	return
		// }

		topTraders := gmgnai.TopTradersToken(chain, tokenAddress)

		for _, trader := range topTraders {

			fmt.Println("Trader: " + trader.Address)

			pnlHistory, err := DeepPNLScanSol(chain, trader.Address, 30)

			if err != nil || pnlHistory == nil {
				//time.Sleep(1 * time.Second)
				continue
			}

			//time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)

			if pnlHistory.SummaryReview.RateBigXPNL > 51 {
				files.AppendToFile("wallet.pnl.txt", trader.Address)
			}

			// time.Sleep(1 * time.Second)
		}

		mongodb.InsertDocumentWithRollback("token_scan_sol", data)
	} else if chain == "eth" {

		tokenAddress = strings.ToLower(tokenAddress)

		filter := bson.M{"contractaddress": tokenAddress}

		_, err := mongodb.FindOne("token_scan_eth", filter)

		if err == nil {
			fmt.Println("Token address already exists in the database.")
			return
		}

		data, err := dexscreener.TokenInfomation("ethereum", tokenAddress)

		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}

		topTraders := gmgnai.TopTradersToken(chain, tokenAddress)

		for _, trader := range topTraders {

			fmt.Println("Trader: " + trader.Address)

			pnlHistory, err := DeepPNLScanETH(chain, trader.Address, 30)

			if err != nil || pnlHistory == nil {
				//time.Sleep(1 * time.Second)
				continue
			}

			//time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)

			if pnlHistory.SummaryReview.RateBigXPNL > 51 {
				files.AppendToFile("wallet.pnl.txt", trader.Address)
			}

			// time.Sleep(1 * time.Second)
		}

		mongodb.InsertDocumentWithRollback("token_scan_eth", data.QI.QuickiAudit)

	} else {
		return
	}
}
