package routers

import (
	"backend-ewallet/controllers"

	"github.com/gin-gonic/gin"
)

func authRouter(r *gin.RouterGroup) {
	authController := controllers.NewAuthController()

	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)
	r.POST("/forgot-password")
	r.POST("/reset-password")
}
