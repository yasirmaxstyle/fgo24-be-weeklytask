package repositories

import (
	"backend-ewallet/models"
	"backend-ewallet/utils"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) CreateTransaction(tx *models.Transaction) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer utils.CloseDB(conn)

	query := `
		INSERT INTO transactions (sender_id, receiver_id, transaction_type, amount, fee,
			description, reference_number, status, created_at, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING transaction_id`

	err = conn.QueryRow(context.Background(), query,
		tx.SenderID, tx.ReceiverID, tx.TransactionType, tx.Amount, tx.Fee,
		tx.Description, tx.ReferenceNumber, tx.Status, tx.CreatedAt,
		tx.Category).Scan(&tx.TransactionID)

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

func (r *TransactionRepository) GetTransactionsByUserID(userID int, limit int) ([]models.Transaction, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return []models.Transaction{}, err
	}
	defer utils.CloseDB(conn)

	query := `
		SELECT transaction_id, sender_id, receiver_id, transaction_type, amount, fee,
			description, reference_number, status, created_at, completed_at, category
		FROM transactions 
		WHERE sender_id = $1 OR receiver_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := conn.Query(context.Background(), query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		err := rows.Scan(
			&tx.TransactionID, &tx.SenderID, &tx.ReceiverID, &tx.TransactionType,
			&tx.Amount, &tx.Fee, &tx.Description, &tx.ReferenceNumber,
			&tx.Status, &tx.CreatedAt, &tx.CompletedAt, &tx.Category)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *TransactionRepository) ProcessTransfer(senderID, receiverID int, amount, fee float64, description, referenceNumber string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer utils.CloseDB(conn)

	// Start transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Get sender balance
	var senderBalance float64
	err = tx.QueryRow(context.Background(), "SELECT balance FROM users WHERE user_id = $1", senderID).Scan(&senderBalance)
	if err != nil {
		return err
	}

	// Check if sender has sufficient balance
	totalAmount := amount + fee
	if senderBalance < totalAmount {
		return pgx.ErrNoRows // Use this to indicate insufficient balance
	}

	// Update sender balance
	newSenderBalance := senderBalance - totalAmount
	_, err = tx.Exec(context.Background(), "UPDATE users SET balance = $1, updated_at = $2 WHERE user_id = $3",
		newSenderBalance, time.Now(), senderID)
	if err != nil {
		return err
	}

	// Update receiver balance
	_, err = tx.Exec(context.Background(), "UPDATE users SET balance = balance + $1, updated_at = $2 WHERE user_id = $3",
		amount, time.Now(), receiverID)
	if err != nil {
		return err
	}

	// Create transaction record
	_, err = tx.Exec(context.Background(), `
		INSERT INTO transactions (sender_id, receiver_id, transaction_type, amount, fee,
			description, reference_number, status, created_at, completed_at, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		senderID, receiverID, "transfer", amount, fee, description,
		referenceNumber, "completed", time.Now(), time.Now(), "transfer")
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit(context.Background())
}
