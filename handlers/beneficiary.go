package handlers

import (
	"math/big"
	"net/http"

	"github.com/Prikshit/fundflow-analysis/models"
	"github.com/Prikshit/fundflow-analysis/services"

	"github.com/gin-gonic/gin"
)

// GetBeneficiaries processes transactions to identify beneficiaries
func GetBeneficiaries(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing address parameter"})
		return
	}

	// Fetch normal transactions
	normalTx, err := services.FetchEtherscanData("account", "txlist", address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	// Fetch internal transactions
	internalTx, err := services.FetchEtherscanData("account", "txlistinternal", address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch internal transactions"})
		return
	}

	// Fetch token transfers
	tokenTx, err := services.FetchEtherscanData("account", "tokentx", address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch token transactions"})
		return
	}

	// Analyze transactions properly
	beneficiaries := analyzeTransactions(normalTx, internalTx, tokenTx)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": beneficiaries})
}

// analyzeTransactions determines the beneficiaries from transaction data
func analyzeTransactions(transactions []models.Transaction, internalTxs []models.Transaction, tokenTransfers []models.Transaction) []gin.H {
	beneficiaryMap := make(map[string]float64)

	// Step 1: Process direct transactions
	for _, tx := range transactions {
		if tx.To != "" {
			beneficiaryMap[tx.To] += parseValue(tx.Value)
		}
	}

	// Step 2: Process internal transactions (Smart Contract Interactions)
	for _, tx := range internalTxs {
		if tx.To != "" {
			beneficiaryMap[tx.To] += parseValue(tx.Value)
		}
	}

	// Step 3: Process Token Transfers (ERC-20, ERC-721, ERC-1155)
	for _, tx := range tokenTransfers {
		if tx.To != "" {
			beneficiaryMap[tx.To] += parseValue(tx.Value)
		}
	}

	// Step 4: Convert to JSON response format
	var results []gin.H
	for address, amount := range beneficiaryMap {
		results = append(results, gin.H{
			"beneficiary_address": address,
			"amount":              amount,
		})
	}

	return results
}

// parseValue converts Wei (string) to Ether (float64)
func parseValue(valueStr string) float64 {
	wei := new(big.Int)
	wei.SetString(valueStr, 10) // Convert string to big.Int

	ether := new(big.Float).SetInt(wei)        // Convert big.Int to big.Float
	divisor := new(big.Float).SetFloat64(1e18) // 1 ETH = 10^18 Wei
	ether.Quo(ether, divisor)                  // Divide Wei by 10^18

	result, _ := ether.Float64() // Convert big.Float to float64
	return result
}
