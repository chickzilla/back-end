package routes

import (
	"github.com/Her_feeling/back-end/utils/middleware"
	"github.com/gin-gonic/gin"
)

func ConfigRouters(server *gin.Engine) {
	server.POST("/result-text", middleware.GetEmailFromToken(), getResultText)
	server.POST("/sign-up", SignUp)
	server.POST("/sign-in", SignIn)
	server.POST("sign-in-sso", SignInWithSSO)
	server.POST("/sign-out", SignOut)
	server.GET("/histories", middleware.GetEmailFromToken(), middleware.AuthMiddleWare(), GetHistories)
}
