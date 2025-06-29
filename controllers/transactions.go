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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get sender user
	sender, err := ctrl.userRepo.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sender information"})
		return
	}

	// Verify PIN
	if !utils.CheckPasswordHash(req.Pin, sender.PinHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid PIN"})
		return
	}

	// Get receiver by phone
	receiver, err := ctrl.userRepo.GetUserByPhone(req.ReceiverPhone)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get receiver information"})
		return
	}

	// Check if sender is trying to transfer to themselves
	if sender.UserID == receiver.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to yourself"})
		return
	}

	// Calculate fee (simple 1% fee)
	fee := req.Amount * 0.01
	totalAmount := req.Amount + fee

	// Check sender balance
	if sender.Balance < totalAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	// Generate reference number
	referenceNumber := utils.GenerateReferenceNumber()

	// Process transfer
	err = ctrl.transactionRepo.ProcessTransfer(
		sender.UserID,
		receiver.UserID,
		req.Amount,
		fee,
		req.Description,
		referenceNumber,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transfer failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer successful",
		"data": gin.H{
			"reference_number": referenceNumber,
			"amount":           req.Amount,
			"fee":              fee,
			"receiver_name":    receiver.FullName,
			"receiver_phone":   receiver.Phone,
		},
	})
}
