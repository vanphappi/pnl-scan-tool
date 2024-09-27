package main

import (
	"fmt"
	"log"
	"os"
	"pnl-solana-tool/package/configs"
	"pnl-solana-tool/platform/database/mongodb"
	"pnl-solana-tool/src/service"
	"strconv"
)

func main() {
	var env, err = configs.LoadConfig(".")

	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	mongoConfig := mongodb.MongoDB{
		DBUsername: env.DBUser,
		DBPassword: env.DBPassword,
		DBHost:     env.DBHost,
		DBPort:     env.DBPort,
		DBName:     env.DBName,
	}

	if err := mongodb.InitMongo(mongoConfig); err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	defer mongodb.Shutdown()

	if len(os.Args) == 2 {
		// Get arguments
		tokenAddress := os.Args[1]

		// Call the PNLScan function with the parsed arguments
		service.TopTraderScan(tokenAddress)
		return
	}

	if len(os.Args) == 3 {

		// Get arguments
		address := os.Args[1]
		numberStr := os.Args[2]

		// Convert number argument to integer
		number, err := strconv.Atoi(numberStr)

		if err != nil {
			fmt.Println("Error: second argument should be a number")
			return
		}

		// Call the PNLScan function with the parsed arguments
		service.PNLScan(address, number)
	}

	if len(os.Args) == 4 && os.Args[1] == "deep" {

		// Get arguments
		address := os.Args[2]
		numberStr := os.Args[3]

		// Convert number argument to integer
		number, err := strconv.Atoi(numberStr)

		if err != nil {
			fmt.Println("Error: second argument should be a number")
			return
		}

		// Call the PNLScan function with the parsed arguments
		service.DeepPNLScan(address, number)
	}

}
