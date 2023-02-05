package main

import (
	employee_controller "HR-system/employee_service/controller"
	employee_storage "HR-system/employee_service/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var server = employee_controller.Server{
	Storage: &employee_storage.Storage{},
	Router:  nil,
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
	jwtKey := os.Getenv("JWT_KEY")

	server.Run(
		host, port, user, dbname, password, "disable", jwtKey,
	)
}
