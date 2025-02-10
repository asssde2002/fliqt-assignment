package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig(file_path string) {
	err := godotenv.Load(file_path)
	if err != nil {
		log.Println("No .env file found")
	}
}

func LoadTimeZone() {
	loc, _ := time.LoadLocation(os.Getenv("TIMEZONE"))
	time.Local = loc
}
