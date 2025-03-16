package models

type HttpError struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}
