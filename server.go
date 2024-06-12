package main

import (
	"github.com/Her_feeling/back-end/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	routes.ConfigRouters(server)
	server.Run(":8080")
}