package dto

type MappedLoyaltyHistoryResponseDTO struct {
	Transactions []TransactionDTO `json:"transactions"`
	Cursor       string           `json:"cursor"`
}
