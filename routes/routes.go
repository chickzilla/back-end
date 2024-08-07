package routes

import "github.com/gin-gonic/gin"

func ConfigRouters(server *gin.Engine) {
	server.GET("/result-text", getResultText)
	server.POST("/sign-up", SignUp)
	server.POST("/sign-in", SignIn)
	server.POST("sign-in-sso", SignInWithSSO)
}
