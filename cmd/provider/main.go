package main

import (
	"correlateme-data-provider/internal/configs"
	"correlateme-data-provider/internal/model"
	"correlateme-data-provider/internal/provider"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	pm, err := provider.InitProviders(cfg)
	if err != nil {
		log.Fatalln("Failed to init providers", err)
	}

	db, err := gorm.Open(postgres.Open(dbConnStr))
	if err != nil {
		log.Fatalln("Failed to connect to DB", err)
	}

	RunProvidersAndSaveObservations(db, cfg, pm)
}

func RunProvidersAndSaveObservations(db *gorm.DB, cfg configs.Config, providers map[string]provider.Provider) {
	for metricID, provider := range providers {
		currentTime := time.Now()

		value, meta, err := provider.Get(currentTime)
		if err != nil {
			fmt.Printf("Error getting value from provider for metric %s: %v\n", metricID, err)
			continue
		}

		observation := model.Observation{
			MetricID:  metricID,
			Value:     value,
			Timestamp: currentTime,
			UserID:    cfg.SystemUserID,
			Meta:      meta,
		}

		result := db.Create(&observation)
		if result.Error != nil {
			fmt.Printf("Error saving observation: %v\n", result.Error)
		}
	}
}
