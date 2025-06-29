package models

type RegisterRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Phone    string `form:"phone" json:"phone" binding:"required"`
	FullName string `form:"full_name" json:"full_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
	Pin      string `form:"pin" json:"pin" binding:"required,len=6"`
}

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type TransferRequest struct {
	ReceiverPhone string  `form:"receiver_phone" json:"receiver_phone" binding:"required"`
	Amount        float64 `form:"amount" json:"amount" binding:"required,gt=0"`
	Description   string  `form:"description" json:"description"`
	Pin           string  `form:"pin" json:"pin" binding:"required,len=6"`
}

type TopUpRequest struct {
	Amount          float64 `form:"amount" json:"amount" binding:"required,gt=0"`
	PaymentMethodID int     `form:"payment_method_id" json:"payment_method_id" binding:"required"`
}
