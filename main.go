package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/pkg/config"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/pkg/database"
)

func main() {
	// Initialize configuration and database
	env := config.LoadENV()

	database.Connection(env)

	// Create a new Gin router
	r := gin.Default()

	fmt.Printf("Starting server on port %s\n", env.APP_PORT)

	// Start the Gin server
	if err := r.Run(env.APP_PORT); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
