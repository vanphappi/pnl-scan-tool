package services

import (
	"fmt"
	"pnl-solana-tool/core/photonsol"
	"pnl-solana-tool/package/files"
	"pnl-solana-tool/platform/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

/*
TopTraderScan scans the Solana blockchain and extracts the top traders of a given token address.
It will check if the token address already exists in the database, and if so, it will return immediately.
It fetches the token information and the top traders of the token, and for each top trader, it will fetch the transactions and calculate the PNL history.
It will then check if the PNL history meets certain conditions, and if so, it will append the wallet address to a file named "wallet.pnl.txt".
It will then insert the token information into the database.
*/
func TopTraderScan(tokenAddress string) {

	filter := bson.M{"tokenaddress": tokenAddress}

	_, err := mongodb.FindOne("token_scan", filter)

	if err == nil {
		fmt.Println("Token address already exists in the database.")
		return
	}

	token := photonsol.Token{
		TokenAddress: tokenAddress,
	}

	data, err := token.TokenInfomation()

	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	topTraders, err := token.TopTraders(data.PoolId)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	for _, trader := range topTraders {

		fmt.Println("Trader: " + trader.Attributes.Signer)

		// if trader.Attributes.BoughtCount < 1 || trader.Attributes.BoughtCount > 3 {
		// 	time.Sleep(1 * time.Second)
		// 	continue
		// }

		// solscan := core.Solscan{
		// 	Address:      trader.Attributes.Signer,
		// 	ExcludeToken: "So11111111111111111111111111111111111111111",
		// 	Flow:         "in",
		// }

		// transactions, err := solscan.GetTransactions(30)

		// if err != nil {
		// 	//files.DeleteFile("wallet.csv")
		// 	time.Sleep(1 * time.Second)

		// 	continue
		// }

		// if len(transactions) > 99 || len(transactions) < 3 {
		// 	//files.DeleteFile("wallet.csv")
		// 	time.Sleep(1 * time.Second)
		// 	continue
		// }

		// time.Sleep(1 * time.Second)

		pnlHistory, err := DeepPNLScan(trader.Attributes.Signer, 30)

		if err != nil || pnlHistory == nil {
			//time.Sleep(1 * time.Second)
			continue
		}

		//time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)

		if pnlHistory.SummaryReview.WinRate > 81.0 || pnlHistory.SummaryReview.RateBigXPNL > 51 {
			files.AppendToFile("wallet.pnl.txt", trader.Attributes.Signer)
		}

		// time.Sleep(1 * time.Second)
	}

	mongodb.InsertDocumentWithRollback("token_scan", data)
}
