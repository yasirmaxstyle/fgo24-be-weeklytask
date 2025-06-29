package routers

import (
	"backend-ewallet/middlewares"

	"github.com/gin-gonic/gin"
)

func CombineRouters(r *gin.Engine) {
	protected := r.Group("/")
	r.Use(middlewares.AuthMiddleware())
	
	authRouter(r.Group("/auth"))
	transactionRouter(protected.Group("/transactions"))
}
