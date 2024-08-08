package middleware

import (
	utils "github.com/Her_feeling/back-end/utils/helper"
	"github.com/gin-gonic/gin"
)

func GetEmailFromToken() gin.HandlerFunc {
	return func(context *gin.Context) {

		token, err := context.Cookie("auth_token")

		if err != nil {
			token = ""
		}

		email := utils.ValidateJWTToken(token)

		context.Set("email", email)

		context.Next()
	}
}
