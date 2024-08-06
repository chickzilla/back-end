package main

import (
	"github.com/Her_feeling/back-end/database"
	"github.com/Her_feeling/back-end/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database.InitDB()
	server := gin.Default()

	server.Use(cors.Default())
	routes.ConfigRouters(server)
	server.Run(":8080")
}
