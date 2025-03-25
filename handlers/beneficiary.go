package handlers

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Prikshit/fundflow-analysis/helpers"
	"github.com/Prikshit/fundflow-analysis/models"
	"github.com/Prikshit/fundflow-analysis/services"

	"github.com/gin-gonic/gin"
)

// GetBeneficiaries processes transactions to identify beneficiaries
func GetBeneficiaries(c *gin.Context) {
	address := strings.ToLower(c.Query("address"))
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing address parameter"})
		return
	}

	var (
		normalTxs   []models.Transaction
		internalTxs []models.Transaction
		tokenTxs    []models.Transaction
		wg          sync.WaitGroup
	)

	wg.Add(3) // Run 3 API calls in parallel

	// Fetch normal transactions
	go func() {
		defer wg.Done()
		var err error
		normalTxs, err = services.FetchEtherscanData("account", "txlist", address)
		if err != nil {
			log.Printf("Error fetching normal transactions: %v", err)
		}
	}()

	// Fetch internal transactions
	go func() {
		defer wg.Done()
		var err error
		internalTxs, err = services.FetchEtherscanData("account", "txlistinternal", address)
		if err != nil {
			log.Printf("Error fetching internal transactions: %v", err)
		}
	}()

	// Fetch token transfers
	go func() {
		defer wg.Done()
		var err error
		tokenTxs, err = services.FetchEtherscanData("account", "tokentx", address)
		if err != nil {
			log.Printf("Error fetching tokentx transactions: %v", err)
		}
	}()

	// Wait for API calls to finish
	wg.Wait()

	// Combine all transactions
	allTxs := append(append(normalTxs, internalTxs...), tokenTxs...)

	beneficiaries := analyzeBeneficiaries(allTxs, address)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": beneficiaries})
}

// analyzeBeneficiaries determines the beneficiaries from transaction data
func analyzeBeneficiaries(transactions []models.Transaction, sourceAddress string) []gin.H {
	beneficiaryMap := make(map[string]float64)
	transactionMap := make(map[string][]gin.H)

	for _, tx := range transactions {
		if tx.To == "" || tx.From == "" {
			continue
		}

		from := strings.ToLower(tx.From)
		to := strings.ToLower(tx.To)

		// Ignore self-transfers and unrelated transactions
		if from == to || from != sourceAddress {
			continue
		}

		amount := helpers.ParseValue(tx.Value)
		if amount == 0 {
			continue
		}

		// Aggregate the amounts per beneficiary
		beneficiaryMap[to] += amount

		// Store transaction details
		transactionMap[to] = append(transactionMap[to], gin.H{
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
