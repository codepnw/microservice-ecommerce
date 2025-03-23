package main

import (
	"log"
	"os"

	"github.com/codepnw/microservice-ecommerce/db"
	"github.com/codepnw/microservice-ecommerce/ecom-api/handler"
	"github.com/codepnw/microservice-ecommerce/ecom-api/server"
	"github.com/codepnw/microservice-ecommerce/ecom-api/store"
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

	st := store.NewMySQLStore(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv, os.Getenv("JWT_SECRET"))

	handler.RegisterRoutes(hdl)
	handler.Start(os.Getenv("APP_PORT"))
}
