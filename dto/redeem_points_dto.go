package dto

type RedeemPointsDTO struct {
	AccountId    string `json:"customer_id"`
	Amount       int    `json:"amount"`
	Description  string `json:"description"`
	RewardTierId string `json:"rewardtier"`
}
