package main

import (
	"correlateme-data-provider/internal/configs"
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

	dbConnStr := os.Getenv("DB")
	configPath := os.Getenv("CONFIG_PATH")
	systemUserUD := os.Getenv("SYSTEM_USER_ID")

	log.Println("Loading config from "+configPath, "")

	cfg, err := configs.ReadConfigs(configPath)
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	cfg.DB = dbConnStr
	cfg.SystemUserID = systemUserUD
	if cfg.Debug {
		log.Printf("Loaded config: %+v", cfg)
	}

}
