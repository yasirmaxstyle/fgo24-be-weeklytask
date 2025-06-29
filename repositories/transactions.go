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

