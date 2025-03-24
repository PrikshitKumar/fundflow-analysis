package main

import (
	"log"
	"os"

	"github.com/Prikshit/fundflow-analysis/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// Register endpoints
	r.GET("/beneficiary", handlers.GetBeneficiaries)
	r.GET("/payer", handlers.GetPayers) // Bonus Task

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port:", port)
	r.Run(":" + port)
}
