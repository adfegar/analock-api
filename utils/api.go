package utils

import (
	"log"
	"net/http"

	"github.com/adfer-dev/analock-api/models"
	"github.com/go-playground/validator/v10"
)

type APIFunc func(res http.ResponseWriter, req *http.Request) error

// Function that parses an APIFunc function to a http.HandlerFunc function
func ParseToHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		if err := f(res, req); err != nil {
			WriteJSON(res, 500, err.Error())
		}

	}
}

func TranslateDbErrorToHttpError(err error) *models.HttpError {
	httpError := &models.HttpError{}

	switch err.(type) {
	case *models.DbNotFoundError:
		httpError.Status = 404
	case *models.DbCouldNotParseItemError:
		httpError.Status = 500
	case *models.DbItemAlreadyExistsError:
		httpError.Status = 400
	default:
		httpError.Status = 500
	}

	httpError.Description = err.Error()

	return httpError
}

func HandleValidation(req *http.Request, body interface{}) []*models.HttpError {
	httpErrors := make([]*models.HttpError, 0)

	if parseErr := ReadJSON(req.Body, body); parseErr != nil {
		log.Println(parseErr)
		if validationErrs, ok := parseErr.(validator.ValidationErrors); ok {
			for _, validationErr := range validationErrs {
				httpErrors = append(httpErrors,
					&models.HttpError{Status: 400, Description: "Field" + validationErr.Field() + " must be provided."})
			}
		} else {
			httpError := models.HttpError{Status: 400, Description: "Not valid JSON."}
			httpErrors = append(httpErrors, &httpError)
		}
	}

	return httpErrors
}
