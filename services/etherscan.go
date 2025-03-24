package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Prikshit/fundflow-analysis/models"
	"github.com/go-resty/resty/v2"
)

const etherscanBaseURL = "https://api.etherscan.io/api"

// fetchEtherscanData makes API calls to fetch transaction data
func FetchEtherscanData(module, action, address string) ([]models.Transaction, error) {
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	client := resty.New().SetTimeout(10 * time.Second)

	url := fmt.Sprintf("%s?module=%s&action=%s&address=%s&apikey=%s", etherscanBaseURL, module, action, address, apiKey)
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status  string               `json:"status"`
		Message string               `json:"message"`
		Result  []models.Transaction `json:"result"`
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.Status != "1" {
		return nil, fmt.Errorf("Etherscan API error: %s", result.Message)
	}

	return result.Result, nil
}
