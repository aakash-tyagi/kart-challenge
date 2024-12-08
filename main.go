package main

import (
	"github.com/aakash-tyagi/kart-challenge/config"
	"github.com/aakash-tyagi/kart-challenge/db"
	"github.com/aakash-tyagi/kart-challenge/server"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()

	// Load configs
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err) // Use logger to log fatal errors
	}

	// Connect to db
	dbClient := db.New(config)

	err = dbClient.Connect()
	if err != nil {
		logger.Fatal(err) // Use logger to log fatal errors
	}

	// Setup the server
	s := server.New(
		dbClient,
		logger,
		config,
	)

	// Start the HTTP server
	s.Start() // Call the Start method to run the server
}
