package dto

type EarnPointsDTO struct {
	AccountId   string `json:"customer_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}
