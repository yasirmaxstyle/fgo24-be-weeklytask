package models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type TransactionResponse struct {
	TransactionID   int     `json:"transaction_id"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	Fee             float64 `json:"fee"`
	ReferenceNumber string  `json:"reference_number"`
	Status          string  `json:"status"`
	NewBalance      float64 `json:"new_balance"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	From        int `json:"from"`
	To          int `json:"to"`
}

type PaginatedResponse struct {
	Data []interface{} `json:"data"`
	Meta Meta          `json:"meta"`
}
