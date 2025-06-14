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
}

var userService services.UserService = &services.UserServiceImpl{}

// @Summary		Get user by ID
// @Description	Get user information by their ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	models.User
// @Failure		404	{object}	models.HttpError
// @Security		BearerAuth
// @Router			/users/{id} [get]
func handleGetUser(res http.ResponseWriter, req *http.Request) error {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])

	user, err := userService.GetUserById(uint(id))

	if err != nil {
		httpErr := utils.TranslateDbErrorToHttpError(err)
		return utils.WriteJSON(res, httpErr.Status, httpErr)
	}

	return utils.WriteJSON(res, 200, user)
}

// @Summary		Get user by email
// @Description	Get user information by their email
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			email	path		string	true	"User email"
// @Success		200		{object}	models.User
// @Failure		404		{object}	models.HttpError
// @Security		BearerAuth
// @Router			/users/{email} [get]
func handleGetUserByEmail(res http.ResponseWriter, req *http.Request) error {
	email := mux.Vars(req)["email"]

	user, err := userService.GetUserByEmail(email)

	if err != nil {
		httpErr := utils.TranslateDbErrorToHttpError(err)
		return utils.WriteJSON(res, httpErr.Status, httpErr)
	}

	return utils.WriteJSON(res, 200, user)
}
