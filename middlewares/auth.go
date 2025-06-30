package middlewares

import (
	"backend-ewallet/models"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC in AuthMiddleware: %v", r)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		godotenv.Load()
		secretKey := os.Getenv("APP_SECRET")
		token := strings.Split(c.GetHeader("Authorization"), "Bearer ")

		if len(token) < 2 {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Unauthorized!",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(token[1])
		rawToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Message: "Token Expired!",
				})
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Token Invalid!",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userIdFloat := rawToken.Claims.(jwt.MapClaims)["user_id"]
		userId := int(userIdFloat.(float64))

		c.Set("user_id", userId)
		c.Next()
	}
}
