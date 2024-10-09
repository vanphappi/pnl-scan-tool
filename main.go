package main

import (
	"fmt"
	"log"
	"os"
	_ "pnl-scan-tool/docs"
	"pnl-scan-tool/package/configs"
	"pnl-scan-tool/platform/database/mongodb"
	"pnl-scan-tool/src/services"
	"strconv"
)

func main() {
	var env, err = configs.LoadConfig(".")

	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	mongoConfig := mongodb.MongoDB{
		DBUsername: env.DB_NAME,
		DBPassword: env.DB_PASSWORD,
		DBHost:     env.DB_HOST,
		DBPort:     env.DB_PORT,
		DBName:     env.DB_NAME,
	}

	if err := mongodb.InitMongo(mongoConfig); err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	defer mongodb.Shutdown()

	if len(os.Args) == 3 && os.Args[1] == "rescan" {
		chain := os.Args[2]
		services.ReScanWalletPNLJob(chain)
	}

	if len(os.Args) == 4 && os.Args[1] == "topholder" {
		chain := os.Args[2]
		tokenAddress := os.Args[3]

		// Call the PNLScan function with the parsed arguments
		services.TopHoldersScan(chain, tokenAddress)
	}

	if len(os.Args) == 4 && os.Args[1] == "toptrader" {
		chain := os.Args[2]
		tokenAddress := os.Args[3]

		// Call the PNLScan function with the parsed arguments
		services.TopTraderScan(chain, tokenAddress)
	}

	if len(os.Args) == 5 && os.Args[1] == "deepscan" {
		chain := os.Args[2]
		// Get arguments
		address := os.Args[3]
		numberStr := os.Args[4]

		// Convert number argument to integer
		number, err := strconv.Atoi(numberStr)

		if err != nil {
			fmt.Println("Error: second argument should be a number")
			return
		}

		if chain == "sol" {
			// Call the PNLScan function with the parsed arguments
			services.DeepPNLScanSol(chain, address, number)

		} else if chain == "eth" {
			// Call the PNLScan function with the parsed arguments
			services.DeepPNLScanETH(chain, address, number)

		} else {
			fmt.Println("Error: chain not supported")
			return
		}
	}

	// if len(os.Args) == 4 && os.Args[1] == "scan" {

	// 	// Get arguments
	// 	address := os.Args[2]
	// 	numberStr := os.Args[3]

	// 	// Convert number argument to integer
	// 	number, err := strconv.Atoi(numberStr)

	// 	if err != nil {
	// 		fmt.Println("Error: second argument should be a number")
	// 		return
	// 	}

	// 	// Call the PNLScan function with the parsed arguments
	// 	services.PNLScan(address, number)
	// }

	// app := fiber.New()

	// // Initialize the Task Manager
	// taskManager := services.NewWalletTrackerTaskManager(4, 10, 2*time.Second)

	// // Serve Swagger UI
	// app.Get("/swagger/*", swagger.HandlerDefault) // Swagger endpoint

	// handlers.WalletTrackerRoutes(app, taskManager)

	// // Start the server
	// app.Listen(":3000")
}
