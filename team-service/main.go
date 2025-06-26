package main

import (
	"log"
	"team-service/config"
	"team-service/database"
	"team-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func main() {
	cfg := config.LoadConfig()
	db := database.Connect(cfg)
	defer database.Close(db)

	router := gin.Default()
	routes.Register(router, db)

	addr := cfg.Host + ":" + cfg.Port
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
