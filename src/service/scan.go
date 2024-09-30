package service

import (
	"fmt"
	"pnl-solana-tool/platform/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

func ReScanWalletPNLJob() {
	collection := "all_time_pnl_wallet"

	filter := bson.M{}

	pnlWalletTracker, err := mongodb.FindDocuments(collection, filter, 0, nil)

	if err != nil {
		fmt.Println(err)
	}

	for _, wallet := range pnlWalletTracker {
		fmt.Println("Scan:", wallet["walletaddress"])
		DeepPNLScan(wallet["walletaddress"].(string), 0)
	}

}
