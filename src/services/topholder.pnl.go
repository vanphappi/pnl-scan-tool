package services

import (
	"fmt"
	gmgnai "pnl-scan-tool/core/gmgn.ai"
	"pnl-scan-tool/package/files"
	"pnl-scan-tool/platform/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

func TopHoldersScan(chain string, tokenAddress string) {

	filter := bson.M{"tokenaddress": tokenAddress, "scantype": "topholders"}

	_, err := mongodb.FindOne("token_scan", filter)

	if err == nil {
		fmt.Println("Token address already exists in the database.")
		return
	}

	topHolers := gmgnai.TopHoldersToken(chain, tokenAddress)

	for _, holder := range topHolers {

		fmt.Println("Holder: " + holder.Address)

		pnlHistory, err := DeepPNLScanSol(chain, holder.Address, 30)

		if err != nil || pnlHistory == nil {
			//time.Sleep(1 * time.Second)
			continue
		}

		//time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)

		if pnlHistory.SummaryReview.RateBigXPNL > 51 {
			files.AppendToFile("wallet.pnl.txt", holder.Address)
		}

		// time.Sleep(1 * time.Second)
	}

	mongodb.InsertDocumentWithRollback("token_scan", map[string]interface{}{
		"tokenaddress": tokenAddress,
		"scantype":     "topholders",
	})
}
