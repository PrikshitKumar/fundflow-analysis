package handlers

import (
	"net/http"

	"github.com/Prikshit/fundflow-analysis/helpers"
	"github.com/Prikshit/fundflow-analysis/models"
	"github.com/Prikshit/fundflow-analysis/services"

	"github.com/gin-gonic/gin"
)

// GetPayers finds sources of incoming funds
func GetPayers(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing address parameter"})
		return
	}

	normalTxs, _ := services.FetchEtherscanData("account", "txlist", address)
	internalTxs, _ := services.FetchEtherscanData("account", "txlistinternal", address)
	tokenTxs, _ := services.FetchEtherscanData("account", "tokentx", address)

	// Combine all transaction types
	allTxs := append(append(normalTxs, internalTxs...), tokenTxs...)

	payers := analyzePayers(allTxs)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": payers})
}

// analyzePayers determines the payers from transaction data
func analyzePayers(transactions []models.Transaction) []gin.H {
	payerMap := make(map[string]float64)

	for _, tx := range transactions {
		if tx.From != "" {
			payerMap[tx.From] += helpers.ParseValue(tx.Value)
		}
	}

	var results []gin.H
	for address, amount := range payerMap {
		results = append(results, gin.H{
			"payer_address": address,
			"amount":        amount,
		})
	}

	return results
}
