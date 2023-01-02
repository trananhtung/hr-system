package main

import (
	"HR-system/employee_api/handlers"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var server = handlers.Server{
	DB:     nil,
	Router: nil,
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	loadEnv()

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DB")
	password := os.Getenv("POSTGRES_PASSWORD")

	server.StartServe(
		host, port, user, dbname, password, "disable",
	)
}
