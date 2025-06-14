package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/adfer-dev/analock-api/constants"
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/services"
	"github.com/adfer-dev/analock-api/utils"
	"github.com/gorilla/mux"
)

var bookRegistrationService services.BookActivityRegistrationService = &services.BookActivityRegistrationServiceImpl{}
var gameRegistrationService services.GameActivityRegistrationService = &services.GameActivityRegistrationServiceImpl{}

func InitActivityRegistrationRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/activityRegistrations/books/user/{id:[0-9]+}", utils.ParseToHandlerFunc(handleGetUserBookActivityRegistrations)).Methods("GET")
	router.HandleFunc("/api/v1/activityRegistrations/games/user/{id:[0-9]+}", utils.ParseToHandlerFunc(handleGetUserGameActivityRegistrations)).Methods("GET")
	router.HandleFunc("/api/v1/activityRegistrations/books", utils.ParseToHandlerFunc(handleCreateBookActivityRegistration)).Methods("POST")
	router.HandleFunc("/api/v1/activityRegistrations/games", utils.ParseToHandlerFunc(handleCreateGameActivityRegistration)).Methods("POST")
}

// @Summary		Get user book activity registrations
// @Description	Get all book activity registrations for a user, optionally filtered by date range
// @Tags			activity registrations
// @Accept			json
// @Produce		json
// @Param			id			path		int	true	"User ID"
// @Param			startDate	query		int	false	"Start date timestamp"
// @Param			endDate		query		int	false	"End date timestamp"
// @Success		200			{array}		models.BookActivityRegistration
// @Failure		400			{object}	models.HttpError
// @Failure		500			{object}	models.HttpError
// @Security		BearerAuth
// @Router			/activityRegistrations/books/user/{id} [get]
func handleGetUserBookActivityRegistrations(res http.ResponseWriter, req *http.Request) error {
	userId, _ := strconv.Atoi(mux.Vars(req)["id"])
	startDateString := req.URL.Query().Get(constants.StartDateQueryParam)
	endDateString := req.URL.Query().Get(constants.EndDateQueryParam)

	if len(startDateString) == 0 || len(endDateString) == 0 {
		userBookRegistrations, err := bookRegistrationService.GetUserBookActivityRegistrations(uint(userId))

		if err != nil {
			return utils.WriteJSON(res, 500, err.Error())
		}

		return utils.WriteJSON(res, 200, userBookRegistrations)
	}

	startDate, startTimeErr := strconv.Atoi(req.URL.Query().Get(constants.StartDateQueryParam))

	if startTimeErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.StartDateQueryParam)})
	}

	endDate, endTimeErr := strconv.Atoi(req.URL.Query().Get(constants.EndDateQueryParam))

	if endTimeErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.EndDateQueryParam)})
	}

	userRegistrations, err := bookRegistrationService.GetUserBookActivityRegistrationsTimeRange(uint(userId), int64(startDate), int64(endDate))

	if err != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: err.Error()})
	}

	return utils.WriteJSON(res, 200, userRegistrations)
}

// @Summary		Get user game activity registrations
// @Description	Get all game activity registrations for a user, optionally filtered by date range
// @Tags			activities
// @Accept			json
// @Produce		json
// @Param			id			path		int	true	"User ID"
// @Param			startDate	query		int	false	"Start date timestamp"
// @Param			endDate		query		int	false	"End date timestamp"
// @Success		200			{array}		models.GameActivityRegistration
// @Failure		400			{object}	models.HttpError
// @Failure		500			{object}	models.HttpError
// @Security		BearerAuth
// @Router			/activityRegistrations/games/user/{id} [get]
func handleGetUserGameActivityRegistrations(res http.ResponseWriter, req *http.Request) error {
	userId, _ := strconv.Atoi(mux.Vars(req)["id"])
	startDateString := req.URL.Query().Get(constants.StartDateQueryParam)
	endDateString := req.URL.Query().Get(constants.EndDateQueryParam)

	if len(startDateString) == 0 || len(endDateString) == 0 {
		userGameRegistrations, err := gameRegistrationService.GetUserGameActivityRegistrations(uint(userId))

		if err != nil {
			return utils.WriteJSON(res, 500, err.Error())
		}

		return utils.WriteJSON(res, 200, userGameRegistrations)
	}

	startDate, startTimeErr := strconv.Atoi(req.URL.Query().Get(constants.StartDateQueryParam))

	if startTimeErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.StartDateQueryParam)})
	}

	endDate, endTimeErr := strconv.Atoi(req.URL.Query().Get(constants.EndDateQueryParam))

	if endTimeErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.EndDateQueryParam)})
	}
	userRegistrations, err := gameRegistrationService.GetUserGameActivityRegistrationsTimeRange(uint(userId), int64(startDate), int64(endDate))

	if err != nil {
		return utils.WriteJSON(res, 400, err.Error())
	}

	return utils.WriteJSON(res, 200, userRegistrations)
}

// @Summary		Create book activity registration
// @Description	Create a new book activity registration
// @Tags			activities
// @Accept			json
// @Produce		json
// @Param			body	body		services.AddBookActivityRegistrationBody	true	"Book activity registration information"
// @Success		200		{object}	models.BookActivityRegistration
// @Failure		400		{object}	models.HttpError
// @Security		BearerAuth
// @Router			/activityRegistrations/books [post]
func handleCreateBookActivityRegistration(res http.ResponseWriter, req *http.Request) error {
	entryBody := services.AddBookActivityRegistrationBody{}

	validationErrs := utils.HandleValidation(req, &entryBody)

	if len(validationErrs) > 0 {
		return utils.WriteJSON(res, 400, validationErrs)
	}

	savedBookRegistration, saveBookRegistrationErr := bookRegistrationService.CreateBookActivityRegistration(&entryBody)

	if saveBookRegistrationErr != nil {
		return utils.WriteJSON(res, 400, saveBookRegistrationErr.Error())
	}

	return utils.WriteJSON(res, 200, savedBookRegistration)
}

// @Summary		Create game activity registration
// @Description	Create a new game activity registration
// @Tags			activities
// @Accept			json
// @Produce		json
// @Param			body	body		services.AddGameActivityRegistrationBody	true	"Game activity registration information"
// @Success		200		{object}	models.GameActivityRegistration
// @Failure		400		{object}	models.HttpError
// @Security		BearerAuth
// @Router			/activityRegistrations/games [post]
func handleCreateGameActivityRegistration(res http.ResponseWriter, req *http.Request) error {
	entryBody := services.AddGameActivityRegistrationBody{}

	validationErrs := utils.HandleValidation(req, &entryBody)

	if len(validationErrs) > 0 {
		return utils.WriteJSON(res, 400, validationErrs)
	}

	savedGameRegistration, saveGameRegistrationErr := gameRegistrationService.CreateGameActivityRegistration(&entryBody)

	if saveGameRegistrationErr != nil {
		return utils.WriteJSON(res, 400, saveGameRegistrationErr.Error())
	}

	return utils.WriteJSON(res, 200, savedGameRegistration)
}
