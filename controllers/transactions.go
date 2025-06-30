package controllers

import (
	"backend-ewallet/models"
	"backend-ewallet/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type TransactionController struct {
	userRepo        *models.UserRepository
	transactionRepo *models.TransactionRepository
}

func NewTransactionController() *TransactionController {
	return &TransactionController{
		userRepo:        models.NewUserRepository(),
		transactionRepo: models.NewTransactionRepository(),
	}
}

func (tc *TransactionController) Transfer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	var req models.TransferRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	sender, err := tc.userRepo.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Sender not found",
		})
		return
	}

	if !utils.CheckPasswordHash(req.Pin, sender.PinHash) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid PIN",
		})
		return
	}

	receiver, err := tc.userRepo.GetUserByPhone(req.ReceiverPhone)
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

	if sender.UserID == receiver.UserID {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Cannot transfer to yourself"})
		return
	}

	fee := req.Amount * 0.01
	totalAmount := req.Amount + fee

	if sender.Balance < totalAmount {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Insufficient balance",
		})
		return
	}

	referenceNumber := utils.GenerateReferenceNumber()

	res, err := tc.transactionRepo.ProcessTransfer(
		sender.UserID,
		receiver.UserID,
		"transfer",
		req.Amount,
		fee,
		req.Description,
		referenceNumber,
		"completed",
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

func (tc *TransactionController) GetTransactionHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	// default limit 10 / page
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	transactions, err := tc.transactionRepo.GetTransactionsByUserID(userID.(int), limit)
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

func (tc *TransactionController) GetBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	user, err := tc.userRepo.GetUserByID(userID.(int))
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

func (tc *TransactionController) Topup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	var req models.TopUpRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := tc.userRepo.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	paymentMethod, err := tc.transactionRepo.GetPaymentMethodByID(req.PaymentMethodID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid payment method",
		})
		return
	}

	if req.Amount < paymentMethod.MinAmount || req.Amount > paymentMethod.MaxAmount {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Amount not within allowed range",
		})
		return
	}

	fee := req.Amount * paymentMethod.FeePercentage / 100

	transaction := &models.Transaction{
		ReceiverID:      &user.UserID,
		TransactionType: "topup",
		PaymentMethodID: &req.PaymentMethodID,
		Amount:          req.Amount,
		Fee:             fee,
		Description:     "Top up via " + paymentMethod.MethodName,
		ReferenceNumber: utils.GenerateReferenceNumber(),
		Status:          "completed",
	}

	if err := tc.transactionRepo.CreateTransaction(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create transaction",
		})
		return
	}

	// Update user balance
	newBalance := user.Balance + req.Amount
	if err := tc.userRepo.UpdateUserBalance(user.UserID, newBalance); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update balance",
		})
		return
	}

	response := models.TransactionResponse{
		TransactionID:   transaction.TransactionID,
		TransactionType: transaction.TransactionType,
		Amount:          req.Amount,
		Fee:             fee,
		ReferenceNumber: transaction.ReferenceNumber,
		Status:          "COMPLETED",
		NewBalance:      newBalance,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Top up successful",
		Data:    response,
	})
}
