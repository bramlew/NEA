package main

import (
	"log"

	"github.com/bramlew/NEA/backend/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env: %v", err)
		return
	}
	handlers.Start()
}
