package services

import (
	"fmt"
	"pnl-scan-tool/platform/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

func ReScanWalletPNLJob(chain string) {
	var collection string

	if chain == "sol" {
		collection = "all_time_pnl_wallet_sol"
		filter := bson.M{}

		pnlWalletTracker, err := mongodb.FindDocuments(collection, filter, 0, nil)

		if err != nil {
			fmt.Println(err)
		}

		for _, wallet := range pnlWalletTracker {
			fmt.Println("Scan:", wallet["walletaddress"])
			DeepPNLScanSol(chain, wallet["walletaddress"].(string), 0)
		}
	} else if chain == "eth" {
		collection := "all_time_pnl_wallet_eth"
		filter := bson.M{}

		pnlWalletTracker, err := mongodb.FindDocuments(collection, filter, 0, nil)

		if err != nil {
			fmt.Println(err)
		}

		for _, wallet := range pnlWalletTracker {
			fmt.Println("Scan:", wallet["walletaddress"])
			DeepPNLScanETH(chain, wallet["walletaddress"].(string), 0)
		}
	}

}
