package models

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Pin      string `json:"pin" binding:"required,len=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TransferRequest struct {
	ReceiverPhone string  `json:"receiver_phone" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Description   string  `json:"description"`
	Pin           string  `json:"pin" binding:"required,len=6"`
}

type TopUpRequest struct {
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethodID int     `json:"payment_method_id" binding:"required"`
}
