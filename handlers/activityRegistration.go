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

func InitActivityRegistrationRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/activityRegistrations/books/user/{userId:[0-9]+}", utils.ParseToHandlerFunc(handleGetUserBookActivityRegistrations)).Methods("GET")
	router.HandleFunc("/api/v1/activityRegistrations/games/user/{userId:[0-9]+}", utils.ParseToHandlerFunc(handleGetUserGameActivityRegistrations)).Methods("GET")
	router.HandleFunc("/api/v1/activityRegistrations/books", utils.ParseToHandlerFunc(handleCreateBookActivityRegistration)).Methods("POST")
	router.HandleFunc("/api/v1/activityRegistrations/games", utils.ParseToHandlerFunc(handleCreateGameActivityRegistration)).Methods("POST")
}

func handleGetUserBookActivityRegistrations(res http.ResponseWriter, req *http.Request) error {
	userId, _ := strconv.Atoi(mux.Vars(req)["userId"])
	startTime, startTimeErr := strconv.Atoi(req.URL.Query().Get(constants.StartDateQueryParam))

	if startTimeErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.StartDateQueryParam)})
	}

	endTime, endTimeErr := strconv.Atoi(req.URL.Query().Get(constants.EndDateQueryParam))

	if endTimeErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.EndDateQueryParam)})
	}

	userRegistrations, err := services.GetUserBookActivityRegistrationsTimeRange(uint(userId), int64(startTime), int64(endTime))

	if err != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: err.Error()})
	}

	return utils.WriteJSON(res, 200, userRegistrations)
}

func handleGetUserGameActivityRegistrations(res http.ResponseWriter, req *http.Request) error {
	userId, _ := strconv.Atoi(mux.Vars(req)["userId"])

	startTime, startTimeErr := strconv.Atoi(req.URL.Query().Get(constants.StartDateQueryParam))

	if startTimeErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.StartDateQueryParam)})
	}

	endTime, endTimeErr := strconv.Atoi(req.URL.Query().Get(constants.EndDateQueryParam))

	if endTimeErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.EndDateQueryParam)})
	}
	userRegistrations, err := services.GetUserGameActivityRegistrationsTimeRange(uint(userId), int64(startTime), int64(endTime))

	if err != nil {
		return utils.WriteJSON(res, 400, err.Error())
	}

	return utils.WriteJSON(res, 200, userRegistrations)
}

func handleCreateBookActivityRegistration(res http.ResponseWriter, req *http.Request) error {
	entryBody := services.AddBookActivityRegistrationBody{}

	validationErrs := utils.HandleValidation(req, &entryBody)

	if len(validationErrs) > 0 {
		return utils.WriteJSON(res, 400, validationErrs)
	}

	savedBookRegistration, saveBookRegistrationErr := services.CreateBookActivityRegistration(&entryBody)

	if saveBookRegistrationErr != nil {
		return utils.WriteJSON(res, 400, saveBookRegistrationErr.Error())
	}

	return utils.WriteJSON(res, 200, savedBookRegistration)
}

func handleCreateGameActivityRegistration(res http.ResponseWriter, req *http.Request) error {
	entryBody := services.AddGameActivityRegistrationBody{}

	validationErrs := utils.HandleValidation(req, &entryBody)

	if len(validationErrs) > 0 {
		return utils.WriteJSON(res, 400, validationErrs)
	}

	savedGameRegistration, saveGameRegistrationErr := services.CreateGameActivityRegistration(&entryBody)

	if saveGameRegistrationErr != nil {
		return utils.WriteJSON(res, 400, saveGameRegistrationErr.Error())
	}

	return utils.WriteJSON(res, 200, savedGameRegistration)
}
