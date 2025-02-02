package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}
}
