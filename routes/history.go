package routes

import (
	"github.com/Her_feeling/back-end/services"
	"github.com/gin-gonic/gin"
)

func GetHistories(context *gin.Context) {
	services.GetUserHistories(context)
}
