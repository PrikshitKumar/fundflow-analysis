package handlers

import (
	"net/http"

	"github.com/Prikshit/fundflow-analysis/helpers"
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
	normalTxs, _ := services.FetchEtherscanData("account", "txlist", address)
	internalTxs, _ := services.FetchEtherscanData("account", "txlistinternal", address)
	tokenTxs, _ := services.FetchEtherscanData("account", "tokentx", address)

	// Combine all transactions
	allTxs := append(append(normalTxs, internalTxs...), tokenTxs...)

	beneficiaries := analyzeBeneficiaries(allTxs)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": beneficiaries})
}

// analyzeBeneficiaries determines the beneficiaries from transaction data
func analyzeBeneficiaries(transactions []models.Transaction) []gin.H {
	beneficiaryMap := make(map[string]float64)
	transactionMap := make(map[string][]gin.H)

	for _, tx := range transactions {
		if tx.To == "" || tx.From == "" {
			continue
		}

		amount := helpers.ParseValue(tx.Value)
		if amount == 0 {
			continue
		}

		// Direct transaction (A -> B)
		beneficiaryMap[tx.To] += amount

		// Store transaction details
		transactionMap[tx.To] = append(transactionMap[tx.To], gin.H{
			"tx_amount":      amount,
			"date_time":      helpers.ParseTimestamp(tx.TimeStamp),
			"transaction_id": tx.Hash,
		})
	}

	// Format response
	var results []gin.H
	for address, amount := range beneficiaryMap {
		results = append(results, gin.H{
			"beneficiary_address": address,
			"amount":              amount,
			"transactions":        transactionMap[address],
		})
	}

	return results
}
