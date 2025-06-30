package models

import (
	"backend-ewallet/utils"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	UserID             int        `json:"user_id" db:"user_id"`
	Email              string     `json:"email" db:"email"`
	Phone              string     `json:"phone" db:"phone"`
	FullName           string     `json:"full_name" db:"full_name"`
	PasswordHash       string     `json:"-" db:"password_hash"`
	PinHash            string     `json:"-" db:"pin_hash"`
	Balance            float64    `json:"balance" db:"balance"`
	RegistrationStatus string     `json:"registration_status" db:"registration_status"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin          *time.Time `json:"last_login" db:"last_login"`
	IsActive           bool       `json:"is_active" db:"is_active"`
}

type Contact struct {
	ContactID     int       `json:"contact_id" db:"contact_id"`
	UserID        int       `json:"user_id" db:"user_id"`
	ContactUserID int       `json:"contact_user_id" db:"contact_user_id"`
	ContactName   string    `json:"contact_name" db:"contact_name"`
	ContactPhone  string    `json:"contact_phone" db:"contact_phone"`
	IsFavorite    bool      `json:"is_favorite" db:"is_favorite"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(user *User) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer utils.CloseDB(conn)

	query := `
		INSERT INTO users (email, phone, full_name, password_hash, pin_hash, balance, 
			registration_status, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING user_id`

	err = conn.QueryRow(context.Background(), query,
		user.Email,
		user.Phone,
		user.FullName,
		user.PasswordHash,
		user.PinHash,
		user.Balance,
		user.RegistrationStatus,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsActive).
		Scan(&user.UserID)

	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return &User{}, err
	}
	defer utils.CloseDB(conn)

	query := `
		SELECT user_id, email, phone, full_name, password_hash, pin_hash, balance,
			registration_status, created_at, updated_at, last_login,
			is_active
		FROM users WHERE email = $1 AND is_active = true`

	row, err := conn.Query(context.Background(), query, email)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow[User](row, pgx.RowToStructByName)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByPhone(phone string) (*User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return &User{}, err
	}
	defer utils.CloseDB(conn)

	query := `
		SELECT user_id, email, phone, full_name, password_hash, pin_hash, balance,
			registration_status, created_at, updated_at, last_login,
			is_active
		FROM users WHERE phone = $1 AND is_active = true`

	row, err := conn.Query(context.Background(), query, phone)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow[User](row, pgx.RowToStructByName)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(userID int) (*User, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return &User{}, err
	}
	defer utils.CloseDB(conn)

	query := `
		SELECT user_id, email, phone, full_name, password_hash, pin_hash, balance,
			registration_status, created_at, updated_at, last_login,
			is_active
		FROM users WHERE user_id = $1 AND is_active = true`

	row, err := conn.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow[User](row, pgx.RowToStructByName)
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

