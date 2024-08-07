package routes

import (
	"github.com/Her_feeling/back-end/services"
	"github.com/gin-gonic/gin"
)

func SignUp(context *gin.Context) {
	services.SignUp(context)
}

func SignIn(context *gin.Context) {
	services.SignIn(context)
}

func SignInWithSSO(context *gin.Context) {
	services.SignInWithSSO(context)
}
