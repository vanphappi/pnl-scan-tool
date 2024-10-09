package services

import (
	"fmt"
	gmgnai "pnl-scan-tool/core/gmgn.ai"
	"pnl-scan-tool/package/utils"
	"pnl-scan-tool/platform/database/mongodb"
	ethmodel "pnl-scan-tool/src/model/eth.model"
	gmaimodel "pnl-scan-tool/src/model/gmai.model"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

func DeepPNLScanETH(chain string, walletAddress string, scanDay int) (*ethmodel.PNL, error) {
	var collection string

	if scanDay == 0 {
		collection = "all_time_pnl_wallet_eth"
	} else {
		collection = "30_day_pnl_wallet_eth"
	}

	var pnlHistory ethmodel.PNL

	pnlHistory.WalletAddress = walletAddress

	filter := bson.M{"walletaddress": walletAddress}

	_, err := mongodb.FindOne(collection, filter)

	if err == nil && scanDay != 0 {
		fmt.Println("Wallet Scan PNL already exists in the database.")

		return nil, err
	}

	transactions := gmgnai.ActivityAllTrade(chain, walletAddress, scanDay)

	totalToken := len(transactions)

	fmt.Println("Scan Total:", strconv.Itoa(totalToken)+" Token")

	count := 0

	fmt.Println("========================================================================================")
	fmt.Println("")

	for _, transaction := range transactions {
		var tradeHistory ethmodel.TradeHistory

		count++

		fmt.Println("Scanning Token Address: " + transaction.TokenAddress)

		if transaction.TokenAddress == "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2" {
			continue
		}

		tradeHistory.TokenAddress = transaction.TokenAddress
		tradeHistory.TokenSymbol = transaction.Token.Symbol //data.TokenSymbol

		fmt.Println("Token Symbol: " + transaction.Token.Symbol)

		fmt.Println("----------------------")

		tradeTransactions := gmgnai.ActivityAllTradeToken(chain, walletAddress, transaction.TokenAddress)

		var countBuy int = 0
		var countSell int = 0
		var CountSellActual int = 0
		var totalTokenBuy float64 = 0
		var totalTokenSell float64 = 0
		var totalTokenSellActual float64 = 0

		var totalETHBuy float64 = 0
		var totalETHSell float64 = 0
		var totalETHSellActual float64 = 0

		var tokenHoldAmount float64 = 0
		var tokenHoldETHAmount float64 = 0

		var priceETHFirstBuy float64 = 0
		var priceETHBestSell float64 = 0

		var lastTracsaction gmaimodel.Activity

		for i := len(tradeTransactions) - 1; i >= 0; i-- {
			tradeTransaction := tradeTransactions[i]

			var eventTrade ethmodel.EventTrade

			if tradeTransaction.QuoteAddress != "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2" {
				continue
			}

			lastTracsaction = tradeTransaction

			eventTrade.EventType = tradeTransaction.EventType
			eventTrade.PriceETH = tradeTransaction.Price
			eventTrade.ETHAmount = utils.ConvertStringToFloat64(tradeTransaction.QuoteAmount)
			eventTrade.TokensAmount = utils.ConvertStringToFloat64(tradeTransaction.TokenAmount)
			eventTrade.Timestamp = tradeTransaction.Timestamp
			eventTrade.DateTime = utils.ConvertTimestampToDate(tradeTransaction.Timestamp)

			tradeHistory.EventTrades = append(tradeHistory.EventTrades, eventTrade)

			fmt.Printf("%#v\n", eventTrade)
			fmt.Println("----------------------")

			if eventTrade.EventType == "buy" {
				if priceETHFirstBuy == 0 {
					priceETHFirstBuy = eventTrade.PriceETH
				}

				totalTokenBuy += eventTrade.TokensAmount

				tokenHoldAmount += eventTrade.TokensAmount

				totalETHBuy += eventTrade.ETHAmount

				countBuy++
			} else if eventTrade.EventType == "sell" {
				if eventTrade.PriceETH > priceETHBestSell {
					priceETHBestSell = eventTrade.PriceETH
				}

				if eventTrade.TokensAmount > tokenHoldAmount {
					totalETHSell += tokenHoldAmount * eventTrade.PriceETH
					totalTokenSell += tokenHoldAmount

					if tokenHoldAmount != 0 {
						countSell++
					}

					tokenHoldAmount = 0

				} else {
					totalETHSell += eventTrade.ETHAmount
					totalTokenSell += eventTrade.TokensAmount

					tokenHoldAmount -= eventTrade.TokensAmount

					countSell++
				}

				totalETHSellActual += eventTrade.ETHAmount
				totalTokenSellActual += eventTrade.TokensAmount

				CountSellActual++
			}
		}

		if len(tradeHistory.EventTrades) == 0 {
			fmt.Println("========================================================================================")
			// time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)
			continue
		}

		tradeHistory.StartTime = tradeHistory.EventTrades[0].DateTime
		tradeHistory.EndTime = tradeHistory.EventTrades[len(tradeHistory.EventTrades)-1].DateTime

		tokenHoldETHAmount = (lastTracsaction.Price * lastTracsaction.Token.Price) / lastTracsaction.PriceUSD * tokenHoldAmount // data.PriceQuote * tokenHoldAmount

		var profitETH float64 = totalETHSell - totalETHBuy + tokenHoldETHAmount
		var profitETHActual float64 = totalETHSellActual - totalETHBuy + tokenHoldETHAmount

		fmt.Println("************************************************************************************")
		fmt.Println("")

		progress := fmt.Sprintf("Scan Progress: %.2f %%", (float64(count)/float64(totalToken))*100.0)
		totalScan := fmt.Sprintf("Total Scan: %d/%d", count, totalToken)

		fmt.Println("Profit ETH: ", profitETH)
		fmt.Println("Profit ETH Actual: ", profitETHActual)

		fmt.Println("")

		fmt.Println(progress)
		fmt.Println(totalScan)

		fmt.Println("")

		pnlHistory.TradeHistory = append(pnlHistory.TradeHistory, tradeHistory)
		pnlHistory.SummaryReview.TotalETHPNLAmount += profitETH
		pnlHistory.SummaryReview.TotalETHPNLAmountActual += profitETHActual

		fmt.Println("PNL ETH: ", pnlHistory.SummaryReview.TotalETHPNLAmount)
		fmt.Println("PNL ETH Actual: ", pnlHistory.SummaryReview.TotalETHPNLAmountActual)

		var xPNL float64 = 0
		var xPNLRate float64 = 0
		var xPNLTrade float64 = 0
		var xPNLRateTrade float64 = 0

		var lostXPNL float64 = 0
		var lostXPNLRate float64 = 0
		var lostXPNLTrade float64 = 0
		var lostXPNLRateTrade float64 = 0

		if profitETH > epsilon {
			pnlHistory.SummaryReview.TotalWin++
			if countBuy != 0 {
				xPNL = (totalETHSell + tokenHoldETHAmount) / totalETHBuy
				xPNLRate = (profitETH / totalETHBuy) * 100
				xPNLTrade = priceETHBestSell / priceETHFirstBuy
				xPNLRateTrade = ((priceETHBestSell - priceETHFirstBuy) / priceETHFirstBuy) * 100
			}
			pnlHistory.XPNLs = append(pnlHistory.XPNLs, ethmodel.XPNL{
				TokenAddress:         tradeHistory.TokenAddress,
				TokenSymbol:          tradeHistory.TokenSymbol,
				CountBuy:             countBuy,
				CountSell:            countSell,
				CountSellActual:      CountSellActual,
				TotalETHBuy:          totalETHBuy,
				TotalETHSell:         totalETHSell,
				TotalETHSellActual:   totalETHSellActual,
				TotalTokenBuy:        totalTokenBuy,
				TotalTokenSell:       totalTokenSell,
				TotalTokenSellActual: totalTokenSellActual,
				TokenHoldAmount:      tokenHoldAmount,
				TokenHoldETHAmount:   tokenHoldETHAmount,
				ProfitETH:            profitETH,
				ProfitETHActual:      profitETHActual,
				XPNL:                 xPNL,
				XPNLRate:             xPNLRate,
				PriceETHFirstBuy:     priceETHFirstBuy,
				PriceETHBestSell:     priceETHBestSell,
				XPNLTrade:            xPNLTrade,
				XPNLRateTrade:        xPNLRateTrade,
				StartTime:            tradeHistory.EventTrades[0].DateTime,
				EndTime:              tradeHistory.EventTrades[len(tradeHistory.EventTrades)-1].DateTime,
			})
		}

		if profitETH <= epsilon && profitETHActual > epsilon {
			pnlHistory.SummaryReview.TotalWin++
			if countBuy != 0 {
				xPNL = (totalETHSell + tokenHoldETHAmount) / totalETHBuy
				xPNLRate = (profitETH / totalETHBuy) * 100
				xPNLTrade = priceETHBestSell / priceETHFirstBuy
				xPNLRateTrade = ((priceETHBestSell - priceETHFirstBuy) / priceETHFirstBuy) * 100
			}
			pnlHistory.XPNLs = append(pnlHistory.XPNLs, ethmodel.XPNL{
				TokenAddress:         tradeHistory.TokenAddress,
				TokenSymbol:          tradeHistory.TokenSymbol,
				CountBuy:             countBuy,
				CountSell:            countSell,
				CountSellActual:      CountSellActual,
				TotalETHBuy:          totalETHBuy,
				TotalETHSell:         totalETHSell,
				TotalETHSellActual:   totalETHSellActual,
				TotalTokenBuy:        totalTokenBuy,
				TotalTokenSell:       totalTokenSell,
				TotalTokenSellActual: totalTokenSellActual,
				TokenHoldAmount:      tokenHoldAmount,
				TokenHoldETHAmount:   tokenHoldETHAmount,
				ProfitETH:            profitETH,
				ProfitETHActual:      profitETHActual,
				XPNL:                 xPNL,
				XPNLRate:             xPNLRate,
				PriceETHFirstBuy:     priceETHFirstBuy,
				PriceETHBestSell:     priceETHBestSell,
				XPNLTrade:            xPNLTrade,
				XPNLRateTrade:        xPNLRateTrade,
				StartTime:            tradeHistory.EventTrades[0].DateTime,
				EndTime:              tradeHistory.EventTrades[len(tradeHistory.EventTrades)-1].DateTime,
			})
		}

		if profitETH <= epsilon && profitETHActual <= epsilon {
			pnlHistory.SummaryReview.TotalLost++
			if countBuy != 0 {
				lostXPNL = (totalETHSell + tokenHoldETHAmount) / totalETHBuy
				lostXPNLRate = (profitETH / totalETHBuy) * 100
				lostXPNLTrade = priceETHBestSell / priceETHFirstBuy
				lostXPNLRateTrade = ((priceETHBestSell - priceETHFirstBuy) / priceETHFirstBuy) * 100
			}
			pnlHistory.LostXPNLs = append(pnlHistory.LostXPNLs, ethmodel.LostXPNL{
				TokenAddress:         tradeHistory.TokenAddress,
				TokenSymbol:          tradeHistory.TokenSymbol,
				CountBuy:             countBuy,
				CountSell:            countSell,
				CountSellActual:      CountSellActual,
				TotalETHBuy:          totalETHBuy,
				TotalETHSell:         totalETHSell,
				TotalETHSellActual:   totalETHSellActual,
				TotalTokenBuy:        totalTokenBuy,
				TotalTokenSell:       totalTokenSell,
				TotalTokenSellActual: totalTokenSellActual,
				TokenHoldAmount:      tokenHoldAmount,
				TokenHoldETHAmount:   tokenHoldETHAmount,
				ProfitETH:            profitETH,
				ProfitETHActual:      profitETHActual,
				LostXPNL:             lostXPNL,
				LostXPNLRate:         lostXPNLRate,
				PriceETHFirstBuy:     priceETHFirstBuy,
				PriceETHBestSell:     priceETHBestSell,
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

		// time.Sleep(time.Duration(generateRandomInt(1000, 2000)) * time.Millisecond)
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
		fmt.Printf(" Total ETH Buy: %.2f ETH  | Total ETH Sell: %.2f ETH | Total ETH Sell Actual: %.2f ETH\n", xpnl.TotalETHBuy, xpnl.TotalETHSell, xpnl.TotalETHSellActual)
		fmt.Printf(" Profit ETH: %.2f ETH | Profit ETH Actual: %.2f ETH\n", xpnl.ProfitETH, xpnl.ProfitETHActual)
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
		fmt.Printf(" Total ETH Buy: %.2f ETH  | Total ETH Sell: %.2f ETH | Total ETH Sell Actual: %.2f ETH\n", xpnl.TotalETHBuy, xpnl.TotalETHSell, xpnl.TotalETHSellActual)
		fmt.Printf(" Profit ETH: %.2f ETH | Profit ETH Actual: %.2f ETH\n", xpnl.ProfitETH, xpnl.ProfitETHActual)
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
	fmt.Println("Total PNL ETH: ", pnlHistory.SummaryReview.TotalETHPNLAmount)
	fmt.Println("Total PNL ETH Actual: ", pnlHistory.SummaryReview.TotalETHPNLAmountActual)
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

	return &pnlHistory, nil

}
