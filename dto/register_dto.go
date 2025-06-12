package dto

type RegisterDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phoneNumber"`
}
