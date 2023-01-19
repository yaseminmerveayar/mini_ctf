package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/migrations"
	"github.com/yaseminmerveayar/mini_ctf/utils"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	database.NewConnection(config)

	migrations.AutoMigrate()

	utils.CreateServer(3000)

}
