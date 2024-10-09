package gmgnai

import (
	"encoding/json"
	"fmt"
	"log"
	gmaimodel "pnl-scan-tool/src/model/gmai.model"
)

// Function to get wallet activities with retry and pagination
func getWalletActivities(chain string, wallet string, cursor string) (*gmaimodel.ApiResponseGMGNAI, error) {
	url := fmt.Sprintf("%s%s?type=buy&type=sell&wallet=%s&limit=%d", baseUrl, chain, wallet, limit)
	fmt.Println(url)
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

func ActivityAllTrade(chain string, wallet string, scanDay int) []gmaimodel.Activity {
	var allActivities []gmaimodel.Activity
	cursor := ""
	// count := 0
	for {
		apiResponse, err := getWalletActivities(chain, wallet, cursor)

		if err != nil {
			log.Fatalf("Error fetching data: %v", err)
		}

		// count += len(apiResponse.Data.Activities)

		// fmt.Println("Scan Token Trade: ", count)

		// Append activities to the slice
		allActivities = append(allActivities, apiResponse.Data.Activities...)

		allActivities = RemoveDuplicates(allActivities)

		fmt.Println("Scan Token Trade: ", len(allActivities))

		if len(allActivities) > scanDay && scanDay != 0 {
			allActivities = allActivities[:scanDay]
			break
		}

		// if len(allActivities) > 300 {
		// 	break
		// }

		// If there's no next cursor, break the loop (end of data)
		if apiResponse.Data.Next == "" {
			break
		}

		cursor = apiResponse.Data.Next
	}

	return RemoveDuplicates(allActivities)
}

func RemoveDuplicates(Activitys []gmaimodel.Activity) []gmaimodel.Activity {
	seen := make(map[string]bool)
	var uniqueTransfers []gmaimodel.Activity

	for _, transfer := range Activitys {
		if !seen[transfer.TokenAddress] {
			seen[transfer.TokenAddress] = true
			uniqueTransfers = append(uniqueTransfers, transfer)
		}
	}

	return uniqueTransfers
}
