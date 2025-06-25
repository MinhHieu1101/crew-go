package main

import (
	"fmt"
	"log"

	"user-service/config"
	"user-service/internal/database"
	"user-service/internal/handlers"
	"user-service/internal/model"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	database.Connect()
	database.DB.AutoMigrate(&model.User{})

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.POST("/users", handlers.Handler)
	r.GET("/users", func(c *gin.Context) {
		c.String(200, "GraphQL endpoint at POST /users")
	})

	addr := fmt.Sprintf("%s:%s", config.GetEnv("HOST", "localhost"), config.GetEnv("PORT", "4000"))
	log.Printf("User service running at http://%s/users", addr)
	log.Fatal(r.Run(addr))
}
