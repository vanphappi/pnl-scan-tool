package gmgnai

import (
	"encoding/json"
	"fmt"
	"log"
)

// Function to get wallet activities with retry and pagination
func getWalletActivities(wallet string, cursor string) (*ApiResponseGMGNAI, error) {
	url := fmt.Sprintf("%s?type=buy&type=sell&wallet=%s&limit=%d", baseUrl, wallet, limit)
	if cursor != "" {
		url += "&cursor=" + cursor
	}

	// Fetch with retry logic
	result, err := fetchWithRetry(url)

	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var apiResponse ApiResponseGMGNAI

	if err := json.Unmarshal(result, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func ActivityAllTrade(wallet string, scanDay int) []Activity {
	var allActivities []Activity
	cursor := ""
	count := 0
	for {
		apiResponse, err := getWalletActivities(wallet, cursor)

		if err != nil {
			log.Fatalf("Error fetching data: %v", err)
		}

		count += len(apiResponse.Data.Activities)

		fmt.Println("Scan Token Trade: ", count)

		// Append activities to the slice
		allActivities = append(allActivities, apiResponse.Data.Activities...)

		allActivities = RemoveDuplicates(allActivities)

		if len(allActivities) > scanDay && scanDay != 0 {
			allActivities = allActivities[:scanDay]
			break
		}

		// If there's no next cursor, break the loop (end of data)
		if apiResponse.Data.Next == "" {
			break
		}

		cursor = apiResponse.Data.Next
	}

	return RemoveDuplicates(allActivities)
}

func RemoveDuplicates(Activitys []Activity) []Activity {
	seen := make(map[string]bool)
	var uniqueTransfers []Activity

	for _, transfer := range Activitys {
		if !seen[transfer.TokenAddress] {
			seen[transfer.TokenAddress] = true
			uniqueTransfers = append(uniqueTransfers, transfer)
		}
	}

	return uniqueTransfers
}
