package main

import (
	"tourist-api/apis"
	database "tourist-api/db"
	"tourist-api/utils"
)

func initialDB() {
	POSTGRES_HOST := utils.ENV_POSTGRES_HOST
	POSTGRES_DB := utils.ENV_POSTGRES_DB
	POSTGRES_USER := utils.ENV_POSTGRES_USER
	POSTGRES_PASSWORD := utils.ENV_POSTGRES_PASSWORD
	POSTGRES_PORT := utils.ENV_POSTGRES_PORT

	database.InitDB(POSTGRES_HOST, POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT)
}

func main() {
	initialDB()
	apis.InitAPI()
}
