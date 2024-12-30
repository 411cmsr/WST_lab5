package main

import (
	"WST_lab5_server/config"

	"WST_lab5_server/internal/database/postgres"
	"WST_lab5_server/internal/httpserver/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	postgres.Init()
	httpServer := gin.Default()

	routes.Init(httpServer)

	httpServer.StaticFile("/favicon.ico", "./favicon.ico")
	err := httpServer.Run(":8095")
	if err != nil {
		fmt.Println(err)
		return
	}
}
