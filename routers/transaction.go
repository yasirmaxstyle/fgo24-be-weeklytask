package routers

import (
	"backend-ewallet/controllers"

	"github.com/gin-gonic/gin"
)

func transactionRouter(r *gin.RouterGroup) {
	transactionController := controllers.NewTransactionController()

	r.POST("/transfer", transactionController.Transfer)
	r.GET("/history", transactionController.GetTransactionHistory)
}
