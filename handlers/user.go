package handlers

import (
	"net/http"
	"strconv"

	"github.com/adfer-dev/analock-api/services"
	"github.com/adfer-dev/analock-api/utils"
	"github.com/gorilla/mux"
)

func InitUserRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/users/{id:[0-9]+}", utils.ParseToHandlerFunc(handleGetUser)).Methods("GET")
	router.HandleFunc("/api/v1/users/{email}", utils.ParseToHandlerFunc(handleGetUserByEmail)).Methods("GET")
	router.HandleFunc("/api/v1/users", utils.ParseToHandlerFunc(handleCreateUser)).Methods("POST")
}

func handleGetUser(res http.ResponseWriter, req *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	user, err := services.GetUserById(uint(id))

	if err != nil {
		httpErr := utils.TranslateDbErrorToHttpError(err)
		return utils.WriteJSON(res, httpErr.Status, httpErr)
	}

	return utils.WriteJSON(res, 200, user)
}

func handleGetUserByEmail(res http.ResponseWriter, req *http.Request) error {
	email := mux.Vars(req)["email"]

	user, err := services.GetUserByEmail(email)

	if err != nil {
		httpErr := utils.TranslateDbErrorToHttpError(err)
		return utils.WriteJSON(res, httpErr.Status, httpErr)
	}

	return utils.WriteJSON(res, 200, user)
}

func handleCreateUser(res http.ResponseWriter, req *http.Request) error {
	userBody := services.UserBody{}

	httpErrors := utils.HandleValidation(req, &userBody)

	if len(httpErrors) > 0 {
		return utils.WriteJSON(res, 403, httpErrors)
	}

	user, err := services.SaveUser(userBody)

	if err != nil {
		httpErr := utils.TranslateDbErrorToHttpError(err)
		return utils.WriteJSON(res, httpErr.Status, httpErr)
	}

	return utils.WriteJSON(res, 201, user)
}
