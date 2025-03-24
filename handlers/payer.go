package handlers

import (
	"net/http"

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

	// Fetch transactions
	normalTx, err := services.FetchEtherscanData("account", "txlist", address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	payers := analyzePayers(normalTx)
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": payers})
}

// analyzePayers determines the payers from transaction data
func analyzePayers(transactions []models.Transaction) []gin.H {
	payerMap := make(map[string]float64)

	for _, tx := range transactions {
		if tx.From != "" {
			payerMap[tx.From] += parseValue(tx.Value)
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
