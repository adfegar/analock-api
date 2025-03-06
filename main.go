package main

import (
	"log"

	"github.com/adfer-dev/analock-api/api"
	"github.com/adfer-dev/analock-api/database"
	"github.com/adfer-dev/analock-api/utils"
	"github.com/joho/godotenv"
)

var logger *utils.CustomLogger = utils.GetCustomLogger()

func main() {
	database.GetDatabaseInstance()
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal(envErr)
	}
	server := api.APIServer{Port: 3000}

	logger.InfoLogger.Printf("Server listening at port %d...\n", server.Port)
	logger.ErrorLogger.Println(server.Run().Error())
}
