package routes

import (
	"net/http"

	"github.com/Her_feeling/back-end/services"
	"github.com/gin-gonic/gin"
)

func getResultText(context *gin.Context){

	prompt := context.Query("prompt")

	result, err := services.SendPrompt(prompt)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, result)

	
}