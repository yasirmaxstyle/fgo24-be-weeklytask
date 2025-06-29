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
	Category        string     `json:"category" db:"category"`
}

type TransferRequest struct {
	ReceiverPhone string  `json:"receiver_phone" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Description   string  `json:"description"`
	Pin           string  `json:"pin" binding:"required,len=6"`
}