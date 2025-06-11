package models

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type LoginRequest struct {
	Phone string `json:"phone"`
}
