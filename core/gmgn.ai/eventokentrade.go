package gmgnai

import (
	"encoding/json"
	"fmt"
	"log"
	gmaimodel "pnl-scan-tool/src/model/gmai.model"
)

// Function to get wallet activities with retry and pagination
func getWalletActivitiesToken(chain string, wallet string, token string, cursor string) (*gmaimodel.ApiResponseGMGNAI, error) {
	url := fmt.Sprintf("%s%s?type=buy&type=sell&wallet=%s&limit=%d&token=%s", baseUrl, chain, wallet, limit, token)
	if cursor != "" {
		url += "&cursor=" + cursor
	}

	// Fetch with retry logic
	result, err := fetchWithRetry(url)

	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var apiResponse gmaimodel.ApiResponseGMGNAI

	if err := json.Unmarshal(result, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func ActivityAllTradeToken(chain string, wallet string, token string) []gmaimodel.Activity {
	var allActivities []gmaimodel.Activity
	cursor := ""
	count := 0
	for {
		apiResponse, err := getWalletActivitiesToken(chain, wallet, token, cursor)

		if err != nil {
			log.Fatalf("Error fetching data: %v", err)
		}

		count += len(apiResponse.Data.Activities)

		if count >= 51 {
			break
		}

		fmt.Println("Scan Total Event Trade: ", count)

		// Append activities to the slice
		allActivities = append(allActivities, apiResponse.Data.Activities...)

		// If there's no next cursor, break the loop (end of data)
		if apiResponse.Data.Next == "" {
			break
		}

		cursor = apiResponse.Data.Next
	}

	return allActivities
}
