package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Prikshit/fundflow-analysis/models"
)

const etherscanBaseURL = "https://api.etherscan.io/api"

// FetchEtherscanData retrieves transaction data
func FetchEtherscanData(module, action, address string) ([]models.Transaction, error) {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	log.Println("🔹 Fetching transactions for address:", address)

	client := &http.Client{Timeout: 60 * time.Second}
	url := fmt.Sprintf("%s?module=%s&action=%s&address=%s&apikey=%s", etherscanBaseURL, module, action, address, apiKey)

	log.Println("🔹 Etherscan API URL:", url)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status  string               `json:"status"`
		Message string               `json:"message"`
		Result  []models.Transaction `json:"result"`
	}

	// Parse JSON response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// Check Etherscan API response status
	if result.Status != "1" {
		return nil, fmt.Errorf("Etherscan API error: %s", result.Message)
	}

	log.Println("✅ Transactions successfully retrieved")
	return result.Result, nil
}
