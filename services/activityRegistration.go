package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

type AddBookActivityRegistrationBody struct {
	InternetArchiveId string `json:"internetArchiveId" validate:"required"`
	RegistrationDate  int64  `json:"registrationDate" validate:"required"`
	UserRefer         uint   `json:"userId" validate:"required"`
}

type AddGameActivityRegistrationBody struct {
	GameName         string `json:"gameName" validate:"required"`
	RegistrationDate int64  `json:"registrationDate" validate:"required"`
	UserRefer        uint   `json:"userId" validate:"required"`
}

var bookActivityRegistrationStorage = &storage.BookActivityRegistrationStorage{}
var gameActivityRegistrationStorage = &storage.GameActivityRegistrationStorage{}

func GetUserBookActivityRegistrations(userId uint) ([]*models.BookActivityRegistration, error) {
	dbUserRegistrations, err := bookActivityRegistrationStorage.GetByUserId(userId)

	if err != nil {
		return nil, err
	}

	return dbUserRegistrations.([]*models.BookActivityRegistration), nil
}

func GetUserGameActivityRegistrations(userId uint) ([]*models.GameActivityRegistration, error) {
	dbUserRegistrations, err := gameActivityRegistrationStorage.GetByUserId(userId)

	if err != nil {
		return nil, err
	}

	return dbUserRegistrations.([]*models.GameActivityRegistration), nil
}

func CreateBookActivityRegistration(addRegistrationBody *AddBookActivityRegistrationBody) (*models.BookActivityRegistration, error) {
	dbActivityRegistration := &models.ActivityRegistration{
		RegistrationDate: addRegistrationBody.RegistrationDate,
		UserRefer:        addRegistrationBody.UserRefer,
	}
	createActivityRegistrationErr := activityRegistrationStorage.Create(dbActivityRegistration)

	if createActivityRegistrationErr != nil {
		return nil, createActivityRegistrationErr
	}

	dbBookActivityRegistration := &models.BookActivityRegistration{
		InternetArchiveIdentifier: addRegistrationBody.InternetArchiveId,
		Registration:              *dbActivityRegistration,
	}

	createBookActivityRegistrationErr := bookActivityRegistrationStorage.Create(dbBookActivityRegistration)

	if createBookActivityRegistrationErr != nil {
		return nil, createBookActivityRegistrationErr
	}

	return dbBookActivityRegistration, nil
}

func CreateGameActivityRegistration(addRegistrationBody *AddGameActivityRegistrationBody) (*models.GameActivityRegistration, error) {

	dbActivityRegistration := &models.ActivityRegistration{
		RegistrationDate: addRegistrationBody.RegistrationDate,
		UserRefer:        addRegistrationBody.UserRefer,
	}
	createActivityRegistrationErr := activityRegistrationStorage.Create(dbActivityRegistration)

	if createActivityRegistrationErr != nil {
		return nil, createActivityRegistrationErr
	}

	dbGameActivityRegistration := &models.GameActivityRegistration{
		GameName:     addRegistrationBody.GameName,
		Registration: *dbActivityRegistration,
	}

	createGameActivityRegistrationErr := gameActivityRegistrationStorage.Create(dbGameActivityRegistration)

	if createGameActivityRegistrationErr != nil {
		return nil, createGameActivityRegistrationErr
	}

	return dbGameActivityRegistration, nil
}
