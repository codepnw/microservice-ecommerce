package main

import (
	"log"
	"os"

	"github.com/codepnw/microservice-ecommerce/db"
	"github.com/joho/godotenv"
)

const envPath = "dev.env"

func main() {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("load .env failed: %v", err)
	}

	db, err := db.NewDatabase(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("error opening database: %v", err)		
	}
	defer db.Close()
	log.Println("successfully connected to database")
}