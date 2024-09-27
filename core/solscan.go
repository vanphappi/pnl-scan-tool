package core

import (
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"pnl-solana-tool/package/files"
	"pnl-solana-tool/package/utils"
	"pnl-solana-tool/platform/database/mongodb"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Define the structure to hold each transaction record.
type Transfer struct {
	Signature    string `json:"signature"`
	Time         string `json:"time"`
	Action       string `json:"action"`
	From         string `json:"from"`
	To           string `json:"to"`
	Amount       string `json:"amount"`
	Decimals     string `json:"decimals"`
	TokenAddress string `json:"tokenAddress"`
}

type Solscan struct {
	Address      string `json:"address"`
	ExcludeToken string `json:"exclude_token"`
	Flow         string `json:"flow"`
}

// GetTransactions fetches the transactions of the given wallet address from Solscan using the Solscan API.
//
// It takes an optional parameter `scanDay` which is the number of days to look back in the transaction history.
// If `scanDay` is 0, it fetches all transactions.
//
// If the wallet address already exists in the database, it will return immediately.
//
// It will then fetch the token information and the top traders of the token, and for each top trader, it will fetch the transactions and calculate the PLN history.
// It will then check if the PLN history meets certain conditions, and if so, it will append the wallet address to a file named "wallet.pnl.txt".
// It will then insert the token information into the database.
//
// The function returns a list of Transfer objects which represent the transactions of the given wallet address.
// If the function fails, it will return an error.
func (s *Solscan) GetTransactions(scanDay int) ([]Transfer, error) {
	url := fmt.Sprintf("https://api-v2.solscan.io/v2/account/transfer/export?address=%s&exclude_token=%s&flow=%s", s.Address, s.ExcludeToken, s.Flow)

	var collection string
	var err error

	if scanDay == 0 {
		collection = "all_time_pnl_wallet"
	} else {
		collection = "30_day_pnl_wallet"
	}

	filter := bson.M{"walletaddress": s.Address}

	_, err = mongodb.FindOne(collection, filter)

	if err == nil && scanDay != 0 {
		fmt.Println("Wallet Scan PNL already exists in the database.")
		return nil, err
	}

	fmt.Println("Requesting URL:", url)

	// Create a new HTTP client with a 30-second timeout
	client := &http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Enable certificate verification
			},
		},
	}

	// Implement retry with exponential backoff
	maxRetries := 3
	var resp *http.Response

	// Retry loop
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Create a new request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set headers to mimic a real browser request
		req.Header.Set("accept", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Referer", "https://solscan.io/") // Set referer to appear like a browser visit
		req.Header.Set("Connection", "keep-alive")

		// Set cookies
		cookie := &http.Cookie{
			Name:  "cf_clearance",
			Value: "bEWP4.WnPCpagmL28BvDpU40873WRenpW5JC7loNEb4-1727067484-1.2.1.1-4WfXUd0zr.qdiYt.VHzwlgJcfHw5YV4G5Gv.ushUj649D8slxqNmMFW9_DWNlsZiLQNZpZ5iYqZS_hTV5PVZz9AG.Rkj8mou7vPKGCvuIsZROoxKTzszvw0NMOgFnCyZ86BVmCOu9YDxpkZaii5QXEsZYTO_FzCuoT9keR5yW4O7So7kSZde5akNgGL5vZfOkGEK5q9pH38zo234CCTIyqOYwveCmTyOMFTTeEHZT4A0fiDKK2bV3EL329tYWlkill1eEY5YfokBrtZ_C7IIeEkkW5SWTpa_PxASAHemcaRwd9zxbTvq32Kz5u947I27YCb_5_xFB2A7IU3plIXewIqAhR8trXDhxN0of28svVQY7xRLj3RUVG1oW7amUtrqGI2dOF0jvIzTksQKDISIh2gmVb9nAMEcMiQG4PNQPWIhTyyXLXpaHWATfHgcQBmG",
			Path:  "/",
		}
		req.AddCookie(cookie)

		// Make the request
		resp, err = client.Do(req)
		if err == nil {
			// Break the loop if the request was successful
			break
		}

		// Print the error and retry after exponential backoff
		fmt.Printf("Attempt %d failed: %v\n", attempt+1, err)
		time.Sleep(time.Second * time.Duration(1<<uint(attempt))) // Exponential backoff
	}

	defer resp.Body.Close()

	fmt.Printf("Response status: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Printf("Response body length: %d bytes\n", len(body))

	// Parse the CSV response
	data := string(body)
	reader := csv.NewReader(strings.NewReader(data))
	records, err := reader.ReadAll()

	// If the response CSV parsing fails, fall back to loading from a file
	if err != nil {

		if !files.FileExists("wallet.csv") {
			_command := fmt.Sprintf("python3 cloudflarebypass/pypass.py %s", url)

			_, _, err = utils.Command(_command)

			if err != nil {
				return nil, err
			}

			utils.RenameCSVFilesToWallet(".")
		}

		// Open and read the CSV file

		fmt.Println("Failed to parse CSV from response, falling back to file")
		file, err := os.Open("wallet.csv")
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err = reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %v", err)
		}
	}

	// Slice to hold the transactions
	var transfers []Transfer

	if len(records) < 2 {
		return nil, fmt.Errorf("insufficient data in CSV")
	}

	// Convert the timestamp to time.Time and adjust for the scan day
	timestamp := utils.StringToInt64(records[1][1])
	timeObj := time.Unix(timestamp, 0)
	adjustedTime := timeObj.AddDate(0, 0, -scanDay)
	timeline := adjustedTime.Unix()

	// Process the CSV records
	for i, record := range records {
		// Skip header row
		if i == 0 {
			continue
		}

		// Skip certain transactions based on conditions
		if record[7] == "So11111111111111111111111111111111111111112" || (utils.StringToInt64(record[1]) < timeline && scanDay != 0) {
			continue
		}

		// Create Transfer object
		transfer := Transfer{
			Signature:    record[0],
			Time:         record[1],
			Action:       record[2],
			From:         record[3],
			To:           record[4],
			Amount:       record[5],
			Decimals:     record[6],
			TokenAddress: record[7],
		}

		// Append the transfer to the list
		transfers = append(transfers, transfer)
	}

	// Remove duplicate transfers
	uniqueTransfers := RemoveDuplicates(transfers)

	return uniqueTransfers, nil
}

// RemoveDuplicates removes duplicates from an array of TokenAddress structs based on the Signature field.
func RemoveDuplicates(transfers []Transfer) []Transfer {
	seen := make(map[string]bool)
	var uniqueTransfers []Transfer

	for _, transfer := range transfers {
		if !seen[transfer.TokenAddress] {
			seen[transfer.TokenAddress] = true
			uniqueTransfers = append(uniqueTransfers, transfer)
		}
	}

	return uniqueTransfers
}

func ReverseTransfers(transfers []Transfer) {
	for i, j := 0, len(transfers)-1; i < j; i, j = i+1, j-1 {
		transfers[i], transfers[j] = transfers[j], transfers[i]
	}
}

func ReadFromFile() string {
	fileData, err := os.ReadFile("datawallet.txt")
	if err != nil {
		return ""
	}

	data := string(fileData)

	return data
}
