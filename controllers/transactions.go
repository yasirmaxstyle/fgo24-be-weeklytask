package controllers

import (
	"backend-ewallet/models"
	"backend-ewallet/repositories"
	"backend-ewallet/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type TransactionController struct {
	userRepo        *repositories.UserRepository
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionController() *TransactionController {
	return &TransactionController{
		userRepo:        repositories.NewUserRepository(),
		transactionRepo: repositories.NewTransactionRepository(),
	}
}

func (ctrl *TransactionController) Transfer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Get sender user
	sender, err := ctrl.userRepo.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Sender not found",
		})
		return
	}

	// Verify PIN
	if !utils.CheckPasswordHash(req.Pin, sender.PinHash) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid PIN",
		})
		return
	}

	// Get receiver by phone
	receiver, err := ctrl.userRepo.GetUserByPhone(req.ReceiverPhone)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Receiver not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to get receiver information"})
		return
	}

	// Check if sender is trying to transfer to themselves
	if sender.UserID == receiver.UserID {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Cannot transfer to yourself"})
		return
	}

	// Calculate fee (simple 1% fee)
	fee := req.Amount * 0.01
	totalAmount := req.Amount + fee

	// Check sender balance
	if sender.Balance < totalAmount {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Insufficient balance",
		})
		return
	}

	// Generate reference number
	referenceNumber := utils.GenerateReferenceNumber()

	// Process transfer
	res, err := ctrl.transactionRepo.ProcessTransfer(
		sender.UserID,
		receiver.UserID,
		"TRANSFER",
		req.Amount,
		fee,
		req.Description,
		referenceNumber,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Insufficient balance",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Transfer failed",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Transfer successful",
		Data:    res,
	})
}

func (ctrl *TransactionController) GetTransactionHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// Get limit from query parameter, default to 10
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	transactions, err := ctrl.transactionRepo.GetTransactionsByUserID(userID.(int), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to get transaction history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction history retrieved successfully",
		"data":    transactions,
	})
}

func (ctrl *TransactionController) GetBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	user, err := ctrl.userRepo.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to get balance",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: false,
		Message: "Balance retrieved successfully",
		Data: gin.H{
			"balance": user.Balance,
		},
	})
}
