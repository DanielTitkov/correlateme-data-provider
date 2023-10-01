package main

import (
	"correlateme-data-provider/internal/configs"
	"correlateme-data-provider/internal/model"
	"correlateme-data-provider/internal/provider"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting data provider")

	startDateFlag := flag.String("start", "", "Start date in YYYY-MM-DD format")
	endDateFlag := flag.String("end", "", "End date in YYYY-MM-DD format")
	metricsFlag := flag.String("metrics", "", "Comma separated list of metric names")

	flag.Parse()

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

	if *startDateFlag != "" && *endDateFlag != "" && *metricsFlag != "" {
		log.Println("Running in backfill mode")

		startDate, err := time.Parse("2006-01-02", *startDateFlag)
		if err != nil {
			log.Fatalf("Error parsing start date: %v\n", err)
		}

		endDate, err := time.Parse("2006-01-02", *endDateFlag)
		if err != nil {
			log.Fatalf("Error parsing end date: %v\n", err)
		}

		metrics := strings.Split(*metricsFlag, "|")

		metricNameToID := make(map[string]string)
		for _, metricConfig := range cfg.Metrics {
			metricNameToID[metricConfig.Name] = metricConfig.ID
		}

		for _, metricName := range metrics {
			log.Println("Processing metric:", metricName)

			metricID, ok := metricNameToID[metricName]
			if !ok {
				fmt.Printf("No metric ID found for metric name: %s\n", metricName)
				continue
			}

			provider, ok := pm[metricID]
			if !ok {
				fmt.Printf("No provider found for metric: %s\n", metricName)
				continue
			}

			for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 0, 1) {
				log.Printf("Processing date %s", d)

				RunProviderForSingleDay(db, cfg, metricID, provider, d)
			}
		}
	} else {
		RunProvidersAndSaveObservations(db, cfg, pm)
	}
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

func RunProviderForSingleDay(db *gorm.DB, cfg configs.Config, metricID string, provider provider.Provider, date time.Time) {
	value, meta, err := provider.Get(date)
	if err != nil {
		fmt.Printf("Error getting value from provider for metric %s on date %v: %v\n", metricID, date, err)
		return
	}

	observation := model.Observation{
		MetricID:  metricID,
		Value:     value,
		Timestamp: date,
		UserID:    cfg.SystemUserID,
		Meta:      meta,
	}

	result := db.Create(&observation)
	if result.Error != nil {
		fmt.Printf("Error saving observation for metric %s on date %v: %v\n", metricID, date, result.Error)
	}
}
