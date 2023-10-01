package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting data provider")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbVariable := os.Getenv("DB")

	fmt.Printf("DB: %s\n", dbVariable)
}
