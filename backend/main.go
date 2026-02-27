package main

import (
	"fmt"
	"log"
	"os"

	d "github.com/ccthomas/board-game/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Board Game Backend...")

	// Load the .env file
	// If args are passed, they will be an array of .env files to load
	// TODO: Come up with better pattern for local development
	// TODO: Do Not deploy this logic to a production environment, without additonal safegaurds
	if len(os.Args) > 1 {
		err := godotenv.Load(os.Args[1:]...)
		if err != nil {
			log.Fatal("Error loading .env file")
			return
		}
	}

	// Only support Postgres at this time
	database := d.NewDatabasePostgres()
	db, err := database.Connect()
	if err != nil {
		fmt.Printf("Error occured connecting to datbase: %s", err)
		return
	}

	migration := d.NewMigrationPostgres()
	err = migration.Up(db)
	if err != nil {
		fmt.Printf("Error occured migrating to datbase: %s", err)
		return
	}

	fmt.Println("Migrated database")
}
