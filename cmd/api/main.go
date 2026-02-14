package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/vizucode/concurent-domain-checker/internal/app"
)

func main() {
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("Error loading .env file")
	}

	app.Run()
}
