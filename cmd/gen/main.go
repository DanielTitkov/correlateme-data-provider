package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting generator")
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/models",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbConnStr := os.Getenv("DB")

	gormdb, _ := gorm.Open(postgres.Open(dbConnStr))
	g.UseDB(gormdb)

	log.Println("Connected to DB")

	g.ApplyBasic(
		// Generate structs from all tables of current database
		g.GenerateAllTable()...,
	)

	// Generate the code
	g.Execute()
}
