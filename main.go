package main

import (
	"go-rest-api/config"
	"go-rest-api/models"
	"go-rest-api/routes"
)

func main() {
	config.SetupTracing()
	config.SetupMetrics()
	config.SetupLogger()
	models.ConnectDatabase()
	router := routes.SetupRouter()
	router.Run(":8084")
}
