package handlers

import (
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

	// Fetch transactions
	normalTx, err := services.FetchEtherscanData("account", "txlist", address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	// Analyze transactions
	beneficiaries := analyzeTransactions(normalTx)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": beneficiaries})
}

// analyzeTransactions determines the beneficiaries from transaction data
func analyzeTransactions(transactions []models.Transaction) []gin.H {
	beneficiaryMap := make(map[string]float64)

	for _, tx := range transactions {
		if tx.To != "" {
			beneficiaryMap[tx.To] += parseValue(tx.Value)
		}
	}

	var results []gin.H
	for address, amount := range beneficiaryMap {
		results = append(results, gin.H{
			"beneficiary_address": address,
			"amount":              amount,
		})
	}

	return results
}

func parseValue(valueStr string) float64 {
	return 0 // Placeholder logic
}
