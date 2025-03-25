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

// GetPayers finds sources of incoming funds
func GetPayers(c *gin.Context) {
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

	payers := analyzePayers(allTxs, address)

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": payers})
}

// analyzePayers determines the payers from transaction data
func analyzePayers(transactions []models.Transaction, targetAddress string) []gin.H {
	payerMap := make(map[string]struct {
		Amount       float64
		Transactions []gin.H
	})

	for _, tx := range transactions {
		if tx.From == "" || tx.To == "" {
			continue
		}

		from := strings.ToLower(tx.From)
		to := strings.ToLower(tx.To)

		// Only consider transactions where the target address is the recipient
		if to != targetAddress {
			continue
		}

		value := helpers.ParseValue(tx.Value)
		payer := payerMap[from]
		payer.Amount += value
		payer.Transactions = append(payer.Transactions, gin.H{
			"tx_amount":      value,
			"date_time":      helpers.ParseTimestamp(tx.TimeStamp),
			"transaction_id": tx.Hash,
		})
		payerMap[from] = payer
	}

	var results []gin.H
	for address, data := range payerMap {
		results = append(results, gin.H{
			"payer_address": address,
			"amount":        data.Amount,
			"transactions":  data.Transactions,
		})
	}

	return results
}
