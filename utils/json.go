package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func WriteJSON(res http.ResponseWriter, status int, value any) error {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)

	return json.NewEncoder(res).Encode(value)
}

func ReadJSON(reader io.Reader, body interface{}) error {
	if deserializeErr := json.NewDecoder(reader).Decode(body); deserializeErr != nil {
		return deserializeErr
	}

	if validationErr := validateBody(body); validationErr != nil {
		return validationErr
	}

	return nil
}

func validateBody(body interface{}) error {
	newValidator := validator.New()

	if err := newValidator.Struct(body); err != nil {
		return err
	}

	return nil
}
