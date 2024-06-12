package routes

import (
	"net/http"

	"github.com/Her_feeling/back-end/services"
	"github.com/gin-gonic/gin"
)

func getResultText(context *gin.Context){

	result, err := services.SendPrompt("So sad bro")

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot get result"})
	}

	context.JSON(http.StatusOK, gin.H{"result": result})
	
}