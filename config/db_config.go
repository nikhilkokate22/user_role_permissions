package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var AppConfig *Config

type Config struct {
	AppPort string

	DBHost           string
	DBPort           string
	DBUser           string
	DBPass           string
	DBName           string
	JwtSecret        string
	JwtExpiryMinutes int
}
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	expiry, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINUTES"))
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRY_MINUTES")
	}

	return &Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPass:           os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
		JwtExpiryMinutes: expiry,
	}
}

