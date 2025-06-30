package routers

import (
	"backend-ewallet/controllers"
	"backend-ewallet/middlewares"

	"github.com/gin-gonic/gin"
)

func transactionRouter(r *gin.RouterGroup) {
	transactionController := controllers.NewTransactionController()
	r.Use(middlewares.AuthMiddleware())

	r.POST("/transfer", transactionController.Transfer)
	r.POST("/topup", transactionController.Topup)
	r.GET("/history", transactionController.GetTransactionHistory)
}
