package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig(file_path string) {
	err := godotenv.Load(file_path)
	if err != nil {
		log.Println("No .env file found")
	}
}
