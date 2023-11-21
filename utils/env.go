package utils

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var ENV_POSTGRES_HOST string
var ENV_POSTGRES_DB string
var ENV_POSTGRES_USER string
var ENV_POSTGRES_PASSWORD string
var ENV_POSTGRES_PORT string
var ENV_FRONTEND_DOMAIN string
var ENV_API_PORT string

func init() {
	ENV_POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	ENV_POSTGRES_DB = os.Getenv("POSTGRES_DB")
	ENV_POSTGRES_USER = os.Getenv("POSTGRES_USER")
	ENV_POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	ENV_POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	ENV_FRONTEND_DOMAIN = os.Getenv("FRONTEND_DOMAIN")
	ENV_API_PORT = os.Getenv("API_PORT")
}
