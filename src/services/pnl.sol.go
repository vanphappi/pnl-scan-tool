package services

import (
	"fmt"
	"pnl-scan-tool/core/photon"
	"pnl-scan-tool/core/solscan"
	"pnl-scan-tool/package/utils"
	"pnl-scan-tool/platform/database/mongodb"
	solmodel "pnl-scan-tool/src/model/sol.model"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/rand"
)

const epsilon = 1e-9

// PNLScan scans the Solana blockchain and returns a PNL struct containing relevant data to be used in the PNL algorithm.
func PNLScan(WalletAddress string, scanDay int) (*solmodel.PNL, error) {

	var collection string

	if scanDay == 0 {
		collection = "all_time_pnl_wallet"
	} else {
		collection = "30_day_pnl_wallet"
	}

	var pnlHistory solmodel.PNL

	pnlHistory.WalletAddress = WalletAddress

	filter := bson.M{"walletaddress": WalletAddress}

	_, err := mongodb.FindOne(collection, filter)

	if err == nil && scanDay != 0 {
		fmt.Println("Wallet Scan PNL already exists in the database.")

		// files.DeleteFile("wallet.csv")

		return nil, err
	}

	solscan := solscan.Solscan{
		Address:      WalletAddress,
		ExcludeToken: "So11111111111111111111111111111111111111111",
		Flow:         "in",
	}

	transactions, err := solscan.GetTransactions(scanDay)

	if err != nil {
		fmt.Println("Error: " + err.Error())

		//files.DeleteFile("wallet.csv")

		return nil, err
	}

	totalToken := len(transactions)

	fmt.Println("Scan Total:", strconv.Itoa(totalToken)+" Token")

	count := 0

	fmt.Println("========================================================================================")
	fmt.Println("")

	for _, transaction := range transactions {
		var tradeHistory solmodel.TradeHistory

		count++

		fmt.Println("Scanning Token Address: " + transaction.TokenAddress)

		if transaction.TokenAddress == "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v" ||
			transaction.TokenAddress == "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB" ||
			transaction.TokenAddress == "JUPyiwrYJFskUPiHa7hkeR8VUtAeFoSYbKedZNsDvCN" {
			fmt.Println("========================================================================================")
			time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)
			continue
		}

		token := photon.Token{
			TokenAddress: transaction.TokenAddress,
		}

		data, err := token.TokenInfomation()

		if err != nil {
			fmt.Println("Error: " + err.Error())
			fmt.Println("========================================================================================")
			time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)
			continue
		}

		tradeHistory.TokenAddress = transaction.TokenAddress
		tradeHistory.TokenSymbol = data.TokenSymbol

		fmt.Println("Token Symbol: " + data.TokenSymbol)

		fmt.Println("----------------------")

		wallet := photon.Wallet{
			WalletAddress: WalletAddress,
		}

		tradeTransactions, err := wallet.Transactions(data.PoolId)

		if err != nil {
			fmt.Println("Error: " + err.Error())
			time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)
			continue
		}

		var countBuy int = 0
		var countSell int = 0
		var CountSellActual int = 0
		var totalTokenBuy float64 = 0
		var totalTokenSell float64 = 0
		var totalTokenSellActual float64 = 0

		var totalSolBuy float64 = 0
		var totalSolSell float64 = 0
		var totalSolSellActual float64 = 0

		var tokenHoldAmount float64 = 0
		var tokenHoldSolAmount float64 = 0

		var priceSolFirstBuy float64 = 0
		var priceSolBestSell float64 = 0

		for i := len(tradeTransactions) - 1; i >= 0; i-- {
			tradeTransaction := tradeTransactions[i]

			var eventTrade solmodel.EventTrade

			eventTrade.Type = tradeTransaction.Attributes.Type
			eventTrade.EventType = tradeTransaction.Attributes.EventType
			eventTrade.PriceSol = utils.ConvertStringToFloat64(tradeTransaction.Attributes.PriceQuote)
			eventTrade.SolAmount = utils.ConvertStringToFloat64(tradeTransaction.Attributes.QuoteAmount)
			eventTrade.TokensAmount = utils.ConvertStringToFloat64(tradeTransaction.Attributes.TokensAmount)
			eventTrade.Timestamp = tradeTransaction.Attributes.Timestamp
			eventTrade.DateTime = utils.ConvertTimestampToDate(tradeTransaction.Attributes.Timestamp)

			tradeHistory.EventTrades = append(tradeHistory.EventTrades, eventTrade)

			fmt.Printf("%#v\n", eventTrade)
			fmt.Println("----------------------")

			if eventTrade.EventType == "create_pool" {
				time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)
				continue
			}

			if eventTrade.Type == "buy" {
				if priceSolFirstBuy == 0 {
					priceSolFirstBuy = eventTrade.PriceSol
				}

				totalTokenBuy += eventTrade.TokensAmount

				tokenHoldAmount += eventTrade.TokensAmount

				totalSolBuy += eventTrade.SolAmount

				countBuy++
			} else if eventTrade.Type == "sell" {
				if eventTrade.PriceSol > priceSolBestSell {
					priceSolBestSell = eventTrade.PriceSol
				}

				if eventTrade.TokensAmount > tokenHoldAmount {
					totalSolSell += tokenHoldAmount * eventTrade.PriceSol
					totalTokenSell += tokenHoldAmount

					if tokenHoldAmount != 0 {
						countSell++
					}

					tokenHoldAmount = 0

				} else {
					totalSolSell += eventTrade.SolAmount
					totalTokenSell += eventTrade.TokensAmount

					tokenHoldAmount -= eventTrade.TokensAmount

					countSell++
				}

				totalSolSellActual += eventTrade.SolAmount
				totalTokenSellActual += eventTrade.TokensAmount

				CountSellActual++
			}
		}

		if len(tradeHistory.EventTrades) == 0 {
			fmt.Println("========================================================================================")
			time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)
			continue
		}

		tradeHistory.StartTime = tradeHistory.EventTrades[0].DateTime
		tradeHistory.EndTime = tradeHistory.EventTrades[len(tradeHistory.EventTrades)-1].DateTime

		tokenHoldSolAmount = data.PriceQuote * tokenHoldAmount

		var profitSol float64 = totalSolSell - totalSolBuy + tokenHoldSolAmount
		var profitSolActual float64 = totalSolSellActual - totalSolBuy + tokenHoldSolAmount

		fmt.Println("************************************************************************************")
		fmt.Println("")

		progress := fmt.Sprintf("Scan Progress: %.2f %%", (float64(count)/float64(totalToken))*100.0)
		totalScan := fmt.Sprintf("Total Scan: %d/%d", count, totalToken)

		fmt.Println("Profit SOL: ", profitSol)
		fmt.Println("Profit SOL Actual: ", profitSolActual)

		fmt.Println("")

		fmt.Println(progress)
		fmt.Println(totalScan)

		fmt.Println("")

		pnlHistory.TradeHistory = append(pnlHistory.TradeHistory, tradeHistory)
		pnlHistory.SummaryReview.TotalSolPNLAmount += profitSol
		pnlHistory.SummaryReview.TotalSolPNLAmountActual += profitSolActual

		fmt.Println("PNL SOL: ", pnlHistory.SummaryReview.TotalSolPNLAmount)
		fmt.Println("PNL SOL Actual: ", pnlHistory.SummaryReview.TotalSolPNLAmountActual)

		var xPNL float64 = 0
		var xPNLRate float64 = 0
		var xPNLTrade float64 = 0
		var xPNLRateTrade float64 = 0

		var lostXPNL float64 = 0
		var lostXPNLRate float64 = 0
		var lostXPNLTrade float64 = 0
		var lostXPNLRateTrade float64 = 0

		if profitSol > epsilon {
			pnlHistory.SummaryReview.TotalWin++
			if countBuy != 0 {
				xPNL = (totalSolSell + tokenHoldSolAmount) / totalSolBuy
				xPNLRate = (profitSol / totalSolBuy) * 100
				xPNLTrade = priceSolBestSell / priceSolFirstBuy
				xPNLRateTrade = ((priceSolBestSell - priceSolFirstBuy) / priceSolFirstBuy) * 100
			}
			pnlHistory.XPNLs = append(pnlHistory.XPNLs, solmodel.XPNL{
				TokenAddress:         tradeHistory.TokenAddress,
				TokenSymbol:          tradeHistory.TokenSymbol,
				CountBuy:             countBuy,
				CountSell:            countSell,
				CountSellActual:      CountSellActual,
				TotalSolBuy:          totalSolBuy,
				TotalSolSell:         totalSolSell,
				TotalSolSellActual:   totalSolSellActual,
				TotalTokenBuy:        totalTokenBuy,
				TotalTokenSell:       totalTokenSell,
				TotalTokenSellActual: totalTokenSellActual,
				TokenHoldAmount:      tokenHoldAmount,
				TokenHoldSolAmount:   tokenHoldSolAmount,
				ProfitSol:            profitSol,
				ProfitSolActual:      profitSolActual,
				XPNL:                 xPNL,
				XPNLRate:             xPNLRate,
				PriceSolFirstBuy:     priceSolFirstBuy,
				PriceSolBestSell:     priceSolBestSell,
				XPNLTrade:            xPNLTrade,
				XPNLRateTrade:        xPNLRateTrade,
				StartTime:            tradeHistory.EventTrades[0].DateTime,
				EndTime:              tradeHistory.EventTrades[len(tradeHistory.EventTrades)-1].DateTime,
			})
		}

		if profitSol <= epsilon && profitSolActual > epsilon {
			pnlHistory.SummaryReview.TotalWin++
			if countBuy != 0 {
				xPNL = (totalSolSell + tokenHoldSolAmount) / totalSolBuy
				xPNLRate = (profitSol / totalSolBuy) * 100
				xPNLTrade = priceSolBestSell / priceSolFirstBuy
				xPNLRateTrade = ((priceSolBestSell - priceSolFirstBuy) / priceSolFirstBuy) * 100
			}
			pnlHistory.XPNLs = append(pnlHistory.XPNLs, solmodel.XPNL{
				TokenAddress:         tradeHistory.TokenAddress,
				TokenSymbol:          tradeHistory.TokenSymbol,
				CountBuy:             countBuy,
				CountSell:            countSell,
				CountSellActual:      CountSellActual,
				TotalSolBuy:          totalSolBuy,
				TotalSolSell:         totalSolSell,
				TotalSolSellActual:   totalSolSellActual,
				TotalTokenBuy:        totalTokenBuy,
				TotalTokenSell:       totalTokenSell,
				TotalTokenSellActual: totalTokenSellActual,
				TokenHoldAmount:      tokenHoldAmount,
				TokenHoldSolAmount:   tokenHoldSolAmount,
				ProfitSol:            profitSol,
				ProfitSolActual:      profitSolActual,
				XPNL:                 xPNL,
				XPNLRate:             xPNLRate,
				PriceSolFirstBuy:     priceSolFirstBuy,
				PriceSolBestSell:     priceSolBestSell,
				XPNLTrade:            xPNLTrade,
				XPNLRateTrade:        xPNLRateTrade,
				StartTime:            tradeHistory.EventTrades[0].DateTime,
				EndTime:              tradeHistory.EventTrades[len(tradeHistory.EventTrades)-1].DateTime,
			})
		}

		if profitSol <= epsilon && profitSolActual <= epsilon {
			pnlHistory.SummaryReview.TotalLost++
			if countBuy != 0 {
				lostXPNL = (totalSolSell + tokenHoldSolAmount) / totalSolBuy
				lostXPNLRate = (profitSol / totalSolBuy) * 100
				lostXPNLTrade = priceSolBestSell / priceSolFirstBuy
				lostXPNLRateTrade = ((priceSolBestSell - priceSolFirstBuy) / priceSolFirstBuy) * 100
			}
			pnlHistory.LostXPNLs = append(pnlHistory.LostXPNLs, solmodel.LostXPNL{
				TokenAddress:         tradeHistory.TokenAddress,
				TokenSymbol:          tradeHistory.TokenSymbol,
				CountBuy:             countBuy,
				CountSell:            countSell,
				CountSellActual:      CountSellActual,
				TotalSolBuy:          totalSolBuy,
				TotalSolSell:         totalSolSell,
				TotalSolSellActual:   totalSolSellActual,
				TotalTokenBuy:        totalTokenBuy,
				TotalTokenSell:       totalTokenSell,
				TotalTokenSellActual: totalTokenSellActual,
				TokenHoldAmount:      tokenHoldAmount,
				TokenHoldSolAmount:   tokenHoldSolAmount,
				ProfitSol:            profitSol,
				ProfitSolActual:      profitSolActual,
				LostXPNL:             lostXPNL,
				LostXPNLRate:         lostXPNLRate,
				PriceSolFirstBuy:     priceSolFirstBuy,
				PriceSolBestSell:     priceSolBestSell,
				LostXPNLTrade:        lostXPNLTrade,
				LostXPNLRateTrade:    lostXPNLRateTrade,
				StartTime:            tradeHistory.EventTrades[0].DateTime,
				EndTime:              tradeHistory.EventTrades[len(tradeHistory.EventTrades)-1].DateTime,
			})
		}

		pnlHistory.SummaryReview.WinRate = (float64(pnlHistory.SummaryReview.TotalWin) / float64(pnlHistory.SummaryReview.TotalWin+pnlHistory.SummaryReview.TotalLost)) * 100.0

		fmt.Println("")

		fmt.Println("************************************************************************************")

		fmt.Println("")

		fmt.Println("========================================================================================")

		fmt.Println("")

		time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)
	}

	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	totalXPLs := len(pnlHistory.XPNLs)
	totalLostXPNLs := len(pnlHistory.LostXPNLs)

	var totalBigXPNL int = 0

	fmt.Println("")
	fmt.Println("WIN ..............................................")
	fmt.Println("")

	for _, xpnl := range pnlHistory.XPNLs {
		if xpnl.XPNL >= 2.0 || xpnl.XPNLTrade >= 2.0 {
			totalBigXPNL++
		}
		fmt.Printf("%s - %.2fx | Trade:%2.fx\n Total Buy: %d | Total Sell: %d | StartTime: %s | EndTime: %s\n", xpnl.TokenSymbol, xpnl.XPNL, xpnl.XPNLTrade, xpnl.CountBuy, xpnl.CountSell, xpnl.StartTime, xpnl.EndTime)
		fmt.Printf(" Total Sol Buy: %.2f SOL  | Total Sol Sell: %.2f SOL | Total Sol Sell Actual: %.2f SOL\n", xpnl.TotalSolBuy, xpnl.TotalSolSell, xpnl.TotalSolSellActual)
		fmt.Printf(" Profit Sol: %.2f SOL | Profit Sol Actual: %.2f SOL\n", xpnl.ProfitSol, xpnl.ProfitSolActual)
		fmt.Printf(" xPNL Rate: %.2f %% | xPNL Rate Trade: %.2f %%\n", xpnl.XPNLRate, xpnl.XPNLRateTrade)
	}

	fmt.Println("")
	fmt.Println("LOST ..............................................")
	fmt.Println("")

	for _, xpnl := range pnlHistory.LostXPNLs {
		if xpnl.LostXPNLTrade >= 2.0 {
			totalBigXPNL++
		}
		fmt.Printf("%s - %.2fx | Trade:%2.fx\n Total Buy: %d | Total Sell: %d | StartTime: %s | EndTime: %s\n", xpnl.TokenSymbol, xpnl.LostXPNL, xpnl.LostXPNLTrade, xpnl.CountBuy, xpnl.CountSell, xpnl.StartTime, xpnl.EndTime)
		fmt.Printf(" Total Sol Buy: %.2f SOL  | Total Sol Sell: %.2f SOL | Total Sol Sell Actual: %.2f SOL\n", xpnl.TotalSolBuy, xpnl.TotalSolSell, xpnl.TotalSolSellActual)
		fmt.Printf(" Profit Sol: %.2f SOL | Profit Sol Actual: %.2f SOL\n", xpnl.ProfitSol, xpnl.ProfitSolActual)
		fmt.Printf(" xPNL Rate: %.2f %% | xPNL Rate Trade: %.2f %%\n", xpnl.LostXPNLRate, xpnl.LostXPNLRateTrade)
	}

	pnlHistory.SummaryReview.BigXPNL = totalBigXPNL
	pnlHistory.SummaryReview.RateBigXPNL = float64(totalBigXPNL) / float64(totalXPLs+totalLostXPNLs) * 100.0

	fmt.Println("")
	fmt.Printf("Big XPNL: %d/%d \n", pnlHistory.SummaryReview.BigXPNL, totalXPLs+totalLostXPNLs)
	fmt.Printf("Rate Big XPNL: %.2f %%\n", pnlHistory.SummaryReview.RateBigXPNL)
	fmt.Println("Total Win: ", pnlHistory.SummaryReview.TotalWin)
	fmt.Println("Total Lost: ", pnlHistory.SummaryReview.TotalLost)
	fmt.Printf("Win Rate: %2.f %%\n", pnlHistory.SummaryReview.WinRate)
	fmt.Println("Total PNL SOL: ", pnlHistory.SummaryReview.TotalSolPNLAmount)
	fmt.Println("Total PNL SOL Actual: ", pnlHistory.SummaryReview.TotalSolPNLAmountActual)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	update := bson.M{"$set": bson.M{
		"tradehistory":  pnlHistory.TradeHistory,
		"xpnls":         pnlHistory.XPNLs,
		"lostxpnls":     pnlHistory.LostXPNLs,
		"summaryreview": pnlHistory.SummaryReview,
	}}

	_, err = mongodb.FindAndUpdateWithRollback(collection, filter, update)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//files.DeleteFile("wallet.csv")

	return &pnlHistory, nil
}

func generateRandomInt(min int, max int) int {
	rand.Seed(uint64(time.Now().UnixNano()))
	return min + rand.Intn(max-min+1)
}
