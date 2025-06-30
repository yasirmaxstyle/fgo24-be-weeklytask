package models

import (
	"backend-ewallet/utils"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	TransactionID   int        `json:"transaction_id" db:"transaction_id"`
	SenderID        *int       `json:"sender_id" db:"sender_id"`
	ReceiverID      *int       `json:"receiver_id" db:"receiver_id"`
	TransactionType string     `json:"transaction_type" db:"transaction_type"`
	PaymentMethodID *int       `json:"payment_method_id" db:"method_id"`
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

type TransactionHistory struct {
	ID                 int
	UserID             int
	TransactionID      int
	TransactionSummary string
	BalanceBefore      float64
	BalanceAfter       float64
	RecordedAt         time.Time
}

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) CreateTransaction(tx *Transaction) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer utils.CloseDB(conn)

	if tx.TransactionType == "transfer" {
		query := `
			INSERT INTO transactions (sender_id, receiver_id, transaction_type, amount, fee,
				description, reference_number, status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING transaction_id`

		err = conn.QueryRow(context.Background(), query,
			tx.SenderID, tx.ReceiverID, tx.TransactionType, tx.Amount, tx.Fee,
			tx.Description, tx.Status, tx.CreatedAt).
			Scan(&tx.TransactionID)
	}

	if tx.TransactionType == "topup" {
		query := `
			INSERT INTO transactions (receiver_id, method_id, transaction_type, amount, fee,
				description, reference_number, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING transaction_id`

		err = conn.QueryRow(context.Background(), query,
			tx.ReceiverID, tx.PaymentMethodID, tx.TransactionType, tx.Amount, tx.Fee,
			tx.Description, tx.ReferenceNumber, tx.Status).
			Scan(&tx.TransactionID)
	}

	fmt.Println(err)
	return err
}

func (r *TransactionRepository) UpdateTransactionStatus(transactionID int, status string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer utils.CloseDB(conn)

	var query string
	var args []interface{}

	if status == "completed" {
		query = `UPDATE transactions SET status = $1, completed_at = $2 WHERE transaction_id = $3`
		args = []interface{}{status, time.Now(), transactionID}
	} else {
		query = `UPDATE transactions SET status = $1 WHERE transaction_id = $2`
		args = []interface{}{status, transactionID}
	}

	_, err = conn.Exec(context.Background(), query, args...)
	return err
}

func (r *TransactionRepository) GetTransactionsByUserID(userID int, limit int) ([]Transaction, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return []Transaction{}, err
	}
	defer utils.CloseDB(conn)

	query := `
	SELECT transaction_id, sender_id, receiver_id, transaction_type, method_id,
	amount, fee, description, reference_number, status, created_at, completed_at
	FROM transactions 
	WHERE sender_id = $1 OR receiver_id = $1
	ORDER BY created_at DESC
	LIMIT $2`

	rows, err := conn.Query(context.Background(), query, userID, limit)
	if err != nil {
		return nil, err
	}

	transactions, err := pgx.CollectRows[Transaction](rows, pgx.RowToStructByName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) ProcessTransfer(senderID int, receiverID int, transactionType string, amount float64, fee float64, description string, referenceNumber string, status string) (*TransactionResponse, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer utils.CloseDB(conn)

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var senderBalance float64
	err = tx.QueryRow(context.Background(), "SELECT balance FROM users WHERE user_id = $1", senderID).Scan(&senderBalance)
	if err != nil {
		return nil, err
	}

	totalAmount := amount + fee
	if senderBalance < totalAmount {
		return nil, pgx.ErrNoRows // Use this to indicate insufficient balance
	}

	newSenderBalance := senderBalance - totalAmount
	_, err = tx.Exec(context.Background(), "UPDATE users SET balance = $1, updated_at = $2 WHERE user_id = $3",
		newSenderBalance, time.Now(), senderID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(context.Background(), "UPDATE users SET balance = balance + $1, updated_at = $2 WHERE user_id = $3",
		amount, time.Now(), receiverID)

	if err != nil {
		return nil, err
	}

	// Create transaction record
	_, err = tx.Exec(context.Background(), `
	INSERT INTO transactions (sender_id, receiver_id, transaction_type, amount, fee,
	description, reference_number, status, created_at, completed_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		senderID, receiverID, transactionType, amount, fee, description,
		referenceNumber, status, time.Now(), time.Now())
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	response := &TransactionResponse{
		ReferenceNumber: referenceNumber,
		TransactionType: transactionType,
		Amount:          amount,
		Fee:             fee,
		Status:          status,
		NewBalance:      senderBalance,
	}

	// Commit transaction
	return response, tx.Commit(context.Background())
}

func (r *TransactionRepository) GetPaymentMethodByID(methodID int) (*PaymentMethod, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer utils.CloseDB(conn)

	query := `
		SELECT method_id, method_name, method_type, is_active, min_amount, max_amount, fee_percentage
		FROM payment_methods WHERE method_id = $1 AND is_active = true`

	row, err := conn.Query(context.Background(), query, methodID)
	if err != nil {
		return nil, err
	}

	method, err := pgx.CollectOneRow[PaymentMethod](row, pgx.RowToStructByName)
	if err != nil {
		return nil, err
	}

	return &method, err
}