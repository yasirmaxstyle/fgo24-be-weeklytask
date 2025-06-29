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
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   "User with this email already exists",
		})
		return
	}

	// Check if phone already exists
	existingUser, _ = ctrl.userRepo.GetUserByPhone(req.Phone)
	if existingUser != nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   "User with this phone number already exists",
		})
		return
	}

	// Hash password and PIN
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to hash password",
		})
		return
	}

	pinHash, err := utils.HashPassword(req.Pin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to hash PIN",
		})
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
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		IsActive:           true,
	}

	if err := ctrl.userRepo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create user",
		})
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.PinHash = ""

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
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
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Database error",
		})
		return
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
		return
	}

	// Update last login
	ctrl.userRepo.UpdateLastLogin(user.UserID)

	// Generate token
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: models.LoginResponse{
			Token: token,
			User:  *user,
		},
	})
}

func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := ctrl.userRepo.GetUserByID(userID.(int))
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Remove sensitive data
	user.PasswordHash = ""
	user.PinHash = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}
