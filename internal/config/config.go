package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	cwd, _ := os.Getwd()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
		log.Println("CWD: ", cwd)
	}
}

var DatabaseURL = os.Getenv("DATABASE_URL")
var EmailUser = os.Getenv("EMAIL_USER")
var EmailPassword = os.Getenv("EMAIL_PASSWORD")
