package api

import (
	"fmt"
	"net/http"

	"github.com/adfer-dev/analock-api/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	Port   int
	router *mux.Router
}

func (server *APIServer) Run() error {
	server.router = mux.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		MaxAge:           86400,
		Debug:            false,
	}).Handler(server.router)

	server.router.Use( /*AuthMiddleware,*/ ValidatePathParams)

	server.initRoutes()

	return http.ListenAndServe(fmt.Sprintf(":%d", server.Port), corsHandler)
}

func (server *APIServer) initRoutes() {
	handlers.InitUserRoutes(server.router)
	handlers.InitAuthRoutes(server.router)
	handlers.InitDiaryEntryRoutes(server.router)
}
