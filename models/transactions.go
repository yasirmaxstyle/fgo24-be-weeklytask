package models

import "time"

type Transaction struct {
	TransactionID   int        `json:"transaction_id" db:"transaction_id"`
	SenderID        int        `json:"sender_id" db:"sender_id"`
	ReceiverID      *int       `json:"receiver_id" db:"receiver_id"`
	TransactionType string     `json:"transaction_type" db:"transaction_type"`
	Amount          float64    `json:"amount" db:"amount"`
	Fee             float64    `json:"fee" db:"fee"`
	Description     string     `json:"description" db:"description"`
	ReferenceNumber string     `json:"reference_number" db:"reference_number"`
	Status          string     `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
}

type PaymentMethod struct {
	MethodID      int     `json:"method_id" db:"method_id"`
	MethodName    string  `json:"method_name" db:"method_name"`
	MethodType    string  `json:"method_type" db:"method_type"`
	IsActive      bool    `json:"is_active" db:"is_active"`
	MinAmount     float64 `json:"min_amount" db:"min_amount"`
	MaxAmount     float64 `json:"max_amount" db:"max_amount"`
	FeePercentage float64 `json:"fee_percentage" db:"fee_percentage"`
}
