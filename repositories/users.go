package repositories

import (
	"backend-ewallet/models"
	"backend-ewallet/utils"
	"context"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer utils.CloseDB(conn)

	query := `
		INSERT INTO users (email, phone, full_name, password_hash, pin_hash, balance, 
			registration_status, is_verified, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING user_id`

	err = conn.QueryRow(context.Background(), query,
		user.Email, user.Phone, user.FullName, user.PasswordHash, user.PinHash,
		user.Balance, user.RegistrationStatus, user.IsVerified, user.CreatedAt,
		user.UpdatedAt, user.IsActive).Scan(&user.UserID)

	return err
}
