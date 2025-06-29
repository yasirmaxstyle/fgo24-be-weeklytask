package repositories

import (
	"backend-ewallet/models"
	"backend-ewallet/utils"
	"context"
	"time"
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

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return &models.User{}, err
	}
	defer utils.CloseDB(conn)

	query := `
		SELECT user_id, email, phone, full_name, password_hash, pin_hash, balance,
			registration_status, is_verified, created_at, updated_at, last_login,
			is_active, email_verification_token, phone_verification_token,
			email_verified_at, phone_verified_at
		FROM users WHERE email = $1 AND is_active = true`

	var user models.User
	err = conn.QueryRow(context.Background(), query, email).Scan(
		&user.UserID, &user.Email, &user.Phone, &user.FullName,
		&user.PasswordHash, &user.PinHash, &user.Balance,
		&user.RegistrationStatus, &user.IsVerified, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.EmailVerificationToken, &user.PhoneVerificationToken,
		&user.EmailVerifiedAt, &user.PhoneVerifiedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByPhone(phone string) (*models.User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return &models.User{}, err
	}
	defer utils.CloseDB(conn)

	query := `
		SELECT user_id, email, phone, full_name, password_hash, pin_hash, balance,
			registration_status, is_verified, created_at, updated_at, last_login,
			is_active, email_verification_token, phone_verification_token,
			email_verified_at, phone_verified_at
		FROM users WHERE phone = $1 AND is_active = true`

	var user models.User
	err = conn.QueryRow(context.Background(), query, phone).Scan(
		&user.UserID, &user.Email, &user.Phone, &user.FullName,
		&user.PasswordHash, &user.PinHash, &user.Balance,
		&user.RegistrationStatus, &user.IsVerified, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.EmailVerificationToken, &user.PhoneVerificationToken,
		&user.EmailVerifiedAt, &user.PhoneVerifiedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(userID int) (*models.User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return &models.User{}, err
	}
	defer utils.CloseDB(conn)

	query := `
		SELECT user_id, email, phone, full_name, password_hash, pin_hash, balance,
			registration_status, is_verified, created_at, updated_at, last_login,
			is_active, email_verification_token, phone_verification_token,
			email_verified_at, phone_verified_at
		FROM users WHERE user_id = $1 AND is_active = true`

	var user models.User
	err = conn.QueryRow(context.Background(), query, userID).Scan(
		&user.UserID, &user.Email, &user.Phone, &user.FullName,
		&user.PasswordHash, &user.PinHash, &user.Balance,
		&user.RegistrationStatus, &user.IsVerified, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLogin, &user.IsActive,
		&user.EmailVerificationToken, &user.PhoneVerificationToken,
		&user.EmailVerifiedAt, &user.PhoneVerifiedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserBalance(userID int, newBalance float64) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer utils.CloseDB(conn)

	query := `UPDATE users SET balance = $1, updated_at = $2 WHERE user_id = $3`
	_, err = conn.Exec(context.Background(), query, newBalance, time.Now(), userID)
	return err
}

func (r *UserRepository) UpdateLastLogin(userID int) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer utils.CloseDB(conn)

	query := `UPDATE users SET last_login = $1 WHERE user_id = $2`
	_, err = conn.Exec(context.Background(), query, time.Now(), userID)
	return err
}
