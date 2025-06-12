package dto

type TransactionDTO struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Points    int    `json:"points"`
	Timestamp string `json:"timestamp"`
}
