package main

import (
	"fmt"
	"os"
	"log"
	"github.com/gin-gonic/gin"
	routes "github.com/MahadevMaahi/gin-jwt-authentication/routes"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Server Starting ...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	
	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success" : "Access granted to api - 1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success" : "Access granted to api - 2"})
	})

	router.Run(":" + port)
	fmt.Println("Server listening on port" + port)
}