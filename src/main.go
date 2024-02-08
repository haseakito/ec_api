package main

import (
	"log"

	"github.com/haseakito/ec_api/routes"
	"github.com/joho/godotenv"
)

//
func init() {
	// Load environment variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}
}

// Entrypoint
func main() {
	// Initialize routers
	e := routes.Init()

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}