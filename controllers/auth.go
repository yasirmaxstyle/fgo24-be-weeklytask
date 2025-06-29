package controllers

import (
	"backend-ewallet/models"
	"backend-ewallet/repositories"
	"backend-ewallet/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthController struct {
	userRepo *repositories.UserRepository
}

func NewAuthController() *AuthController {
	return &AuthController{
		userRepo: repositories.NewUserRepository(),
	}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	existingUser, _ := ctrl.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Check if phone already exists
	existingUser, _ = ctrl.userRepo.GetUserByPhone(req.Phone)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this phone already exists"})
		return
	}

	// Hash password and PIN
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	pinHash, err := utils.HashPassword(req.Pin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash PIN"})
		return
	}

	// Create user
	user := &models.User{
		Email:              req.Email,
		Phone:              req.Phone,
		FullName:           req.FullName,
		PasswordHash:       passwordHash,
		PinHash:            pinHash,
		Balance:            0.0,
		RegistrationStatus: "completed",
		IsVerified:         false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		IsActive:           true,
	}

	if err := ctrl.userRepo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.PinHash = ""

	response := models.AuthResponse{
		User:  *user,
		Token: token,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := ctrl.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update last login
	ctrl.userRepo.UpdateLastLogin(user.UserID)

	// Generate token
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.PinHash = ""

	response := models.AuthResponse{
		User:  *user,
		Token: token,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}