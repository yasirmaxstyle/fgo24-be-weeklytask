package models

import (
	"time"
)

type User struct {
	UserID                 int        `json:"user_id" db:"user_id"`
	Email                  string     `json:"email" db:"email"`
	Phone                  string     `json:"phone" db:"phone"`
	FullName               string     `json:"full_name" db:"full_name"`
	PasswordHash           string     `json:"-" db:"password_hash"`
	PinHash                string     `json:"-" db:"pin_hash"`
	Balance                float64    `json:"balance" db:"balance"`
	RegistrationStatus     string     `json:"registration_status" db:"registration_status"`
	IsVerified             bool       `json:"is_verified" db:"is_verified"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin              *time.Time `json:"last_login" db:"last_login"`
	IsActive               bool       `json:"is_active" db:"is_active"`
	EmailVerificationToken string     `json:"-" db:"email_verification_token"`
	PhoneVerificationToken string     `json:"-" db:"phone_verification_token"`
	EmailVerifiedAt        *time.Time `json:"email_verified_at" db:"email_verified_at"`
	PhoneVerifiedAt        *time.Time `json:"phone_verified_at" db:"phone_verified_at"`
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

type UserSession struct {
	SessionID    int       `json:"session_id" db:"session_id"`
	UserID       int       `json:"user_id" db:"user_id"`
	SessionToken string    `json:"session_token" db:"session_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	DeviceInfo   string    `json:"device_info" db:"device_info"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
}