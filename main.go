package main

import (
	"log"
	"todo-api/config"
	"todo-api/database"
	"todo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.InitConfig()

	database.InitDB(cfg)

	r := gin.Default()

	routes.SetupRoutes(r)

	log.Printf("Server starting on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
