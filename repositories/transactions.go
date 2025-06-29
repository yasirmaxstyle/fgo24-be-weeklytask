package repositories

import (
	"backend-ewallet/models"
	"backend-ewallet/utils"
	"context"
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
