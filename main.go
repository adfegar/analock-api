package main

import (
	"log"

	"github.com/adfer-dev/analock-api/api"
	"github.com/adfer-dev/analock-api/utils"
	"github.com/joho/godotenv"
)

//	@title			Analock API
//	@version		1.0
//	@description	This is the API server for the Analock application.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/api/v1
//	@schemes					http https
//	@securityDefinitions.apiKey	Bearer Token
//	@in							header
//	@name						Authorization

var logger *utils.CustomLogger = utils.GetCustomLogger()

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("No env file is present")
	}

	server := api.APIServer{Port: 3000}

	logger.InfoLogger.Printf("Server listening at port %d...\n", server.Port)
	logger.ErrorLogger.Println(server.Run().Error())
}
