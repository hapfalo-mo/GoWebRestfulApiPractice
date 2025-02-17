package main

import (
	"fmt"
	"log"
	"my-gin-app/db"
	"my-gin-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	db.ConnectDB()
	defer db.DB.Close()

	// Intialize Gin Router
	router := gin.Default()

	// Register User API Routes
	routes.SetUsertRoutes(router)

	// Start the Server
	port := ":8080"
	fmt.Println("🚀 Server running on http://localhost" + port)
	if err := router.Run(port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
