package middleware

import (
	"fmt"
	"net/http"

	"github.com/Her_feeling/back-end/database"
	utils "github.com/Her_feeling/back-end/utils/helper"
	"github.com/gin-gonic/gin"
)

func GetEmailFromToken() gin.HandlerFunc {
	return func(context *gin.Context) {

		token, err := context.Cookie("auth_token")

		if err != nil {
			token = ""
		}

		fmt.Println("token", token)

		email := utils.ValidateJWTToken(token)

		context.Set("email", email)

		context.Next()
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		var DB = database.DB
		email, ok := context.Get("email")

		if !ok {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			context.Abort()
		}

		var userId int
		err := DB.QueryRow("SELECT id FROM user WHERE email = ?", email).Scan(&userId)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "email don't exist"})
			context.Abort()
		}

		context.Set("userId", userId)

		context.Next()
	}
}
