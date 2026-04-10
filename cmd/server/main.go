package main

import (
	"log"

	"marketplace-api/config"
	"marketplace-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDB()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	routes.SetupRoutes(r, db)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
