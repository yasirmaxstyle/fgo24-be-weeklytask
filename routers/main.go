package routers

import (
	"github.com/gin-gonic/gin"
)

func CombineRouters(r *gin.Engine) {
	authRouter(r.Group("/auth"))
	transactionRouter(r.Group("/transactions"))
}
