package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	CTX = context.Background()
	RDB *redis.Client
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       getEnvAsInt("REDIS_DB", 0),
	})

	_, err := RDB.Ping(CTX).Result()
	if err != nil {
		log.Fatal("Failed to connect to redis:", err)
	}

}

func CloseRedis() {
	if err := RDB.Close(); err != nil {
		log.Println("Failed to close redis:", err)
		return
	}
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}
