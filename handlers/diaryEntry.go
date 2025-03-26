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

func InitDiaryEntryRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/diaryEntries/user/{userId:[0-9]+}", utils.ParseToHandlerFunc(handleGetUserEntries)).Methods("GET")
	router.HandleFunc("/api/v1/diaryEntries", utils.ParseToHandlerFunc(handleCreateDiaryEntry)).Methods("POST")
	router.HandleFunc("/api/v1/diaryEntries/{id:[0-9]+}", utils.ParseToHandlerFunc(handleUpdateDiaryEntry)).Methods("PUT")
}

func handleGetUserEntries(res http.ResponseWriter, req *http.Request) error {
	userId, _ := strconv.Atoi(mux.Vars(req)["userId"])

	startDateString := req.URL.Query().Get(constants.StartDateQueryParam)
	endDateString := req.URL.Query().Get(constants.EndDateQueryParam)

	if len(startDateString) == 0 || len(endDateString) == 0 {
		userDiaryEntries, err := services.GetUserEntries(uint(userId))

		if err != nil {
			return utils.WriteJSON(res, 500, err.Error())
		}

		return utils.WriteJSON(res, 200, userDiaryEntries)
	}

	startDate, startDateErr := strconv.Atoi(startDateString)

	if startDateErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.StartDateQueryParam)})
	}
	endDate, endDateErr := strconv.Atoi(endDateString)

	if endDateErr != nil {
		return utils.WriteJSON(res, 400, models.HttpError{Status: http.StatusBadRequest, Description: fmt.Sprintf(constants.QueryParamError, constants.EndDateQueryParam)})
	}
	dateIntervalUserDiaryEntries, err := services.GetUserEntriesTimeRange(uint(userId), int64(startDate), int64(endDate))
	if err != nil {
		return utils.WriteJSON(res, 500, err.Error())
	}

	return utils.WriteJSON(res, 200, dateIntervalUserDiaryEntries)
}

func handleCreateDiaryEntry(res http.ResponseWriter, req *http.Request) error {
	entryBody := services.SaveDiaryEntryBody{}

	validationErrs := utils.HandleValidation(req, &entryBody)

	if len(validationErrs) > 0 {
		return utils.WriteJSON(res, 400, validationErrs)
	}

	savedEntry, saveEntryErr := services.SaveDiaryEntry(&entryBody)

	if saveEntryErr != nil {
		return utils.WriteJSON(res, 500, saveEntryErr.Error())
	}

	return utils.WriteJSON(res, 201, savedEntry)
}

func handleUpdateDiaryEntry(res http.ResponseWriter, req *http.Request) error {
	entryId, _ := strconv.Atoi(mux.Vars(req)["id"])
	updateEntryBody := services.UpdateDiaryEntryBody{}

	validationErrs := utils.HandleValidation(req, &updateEntryBody)

	if len(validationErrs) > 0 {
		return utils.WriteJSON(res, 400, validationErrs)
	}

	updatedEntry, updateEntryErr := services.UpdateDiaryEntry(uint(entryId), &updateEntryBody)

	if updateEntryErr != nil {
		return utils.WriteJSON(res, 500, updateEntryErr.Error())
	}

	return utils.WriteJSON(res, 200, updatedEntry)
}
