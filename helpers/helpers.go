package helpers

import (
	"math/big"
	"strconv"
	"time"
)

func ParseValue(valueStr string) float64 {
	wei := new(big.Int)
	wei.SetString(valueStr, 10)

	ether := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetFloat64(1e18)
	ether.Quo(ether, divisor)

	result, _ := ether.Float64()
	return result
}

func ParseTimestamp(timestamp string) string {
	seconds, _ := strconv.ParseInt(timestamp, 10, 64)
	t := time.Unix(seconds, 0)
	return t.Format("2006-01-02 15:04:05")
}

// Helper function to recursively find the final recipient
func TraceFinalBeneficiary(address string, nextHop map[string]string) string {
	seen := make(map[string]bool)
	for {
		next, exists := nextHop[address]
		if !exists || seen[next] {
			break
		}
		seen[next] = true
		address = next
	}
	return address
}
