package main

import (
	"fmt"
	"strings"

	"github.com/Her_feeling/back-end/database"
	"github.com/Her_feeling/back-end/routes"
	utils "github.com/Her_feeling/back-end/utils/helper"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database.InitDB()
	server := gin.Default()
	fmt.Println("Server already deploy with cloud build")

	whiteList, _ := utils.GetEnvNoCon("WHITE_LIST")

	whiteList = strings.TrimSuffix(whiteList, ",")

	server.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(whiteList, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	port, _ := utils.GetEnvNoCon("PORT")

	routes.ConfigRouters(server)
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}
