package models

// Transaction represents an Ethereum transaction
type Transaction struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	Hash        string `json:"hash"`
	TimeStamp   string `json:"timeStamp"`
	TokenSymbol string `json:"tokenSymbol"`
}
