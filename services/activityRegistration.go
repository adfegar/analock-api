package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

// Storage interfaces
type BookActivityRegistrationStorageInterface interface {
	GetByUserId(userId uint) (interface{}, error)
	GetByUserIdAndTimeRange(userId uint, startTime int64, endTime int64) (interface{}, error)
	Create(data interface{}) error
}

type GameActivityRegistrationStorageInterface interface {
	GetByUserId(userId uint) (interface{}, error)
	GetByUserIdAndInterval(userId uint, startDate int64, endDate int64) (interface{}, error)
	Create(data interface{}) error
}

type ActivityRegistrationStorageInterface interface {
	Create(data interface{}) error
	Update(data interface{}) error
	Delete(id uint) error
}

// BookActicityRegistrationService interface and implementation
type BookActivityRegistrationService interface {
	GetUserBookActivityRegistrations(userId uint) ([]*models.BookActivityRegistration, error)
	GetUserBookActivityRegistrationsTimeRange(userId uint, startTime int64, endTime int64) ([]*models.BookActivityRegistration, error)
	CreateBookActivityRegistration(addRegistrationBody *AddBookActivityRegistrationBody) (*models.BookActivityRegistration, error)
}

type DefaultBookActivityRegistrationService struct{}

// GameActicityRegistrationService interface and implementation
type GameActivityRegistrationService interface {
	GetUserGameActivityRegistrations(userId uint) ([]*models.GameActivityRegistration, error)
	GetUserGameActivityRegistrationsTimeRange(userId uint, startDate int64, endDate int64) ([]*models.GameActivityRegistration, error)
	CreateGameActivityRegistration(addRegistrationBody *AddGameActivityRegistrationBody) (*models.GameActivityRegistration, error)
}

type DefaultGameActivityRegistrationService struct{}

// Request bodies structs
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

var bookActivityRegistrationStorage BookActivityRegistrationStorageInterface = &storage.BookActivityRegistrationStorage{}
var gameActivityRegistrationStorage GameActivityRegistrationStorageInterface = &storage.GameActivityRegistrationStorage{}
var activityRegistrationStorage ActivityRegistrationStorageInterface = &storage.ActivityRegistrationStorage{}

func (defaultBookActivityRegistrationService *DefaultBookActivityRegistrationService) GetUserBookActivityRegistrations(userId uint) ([]*models.BookActivityRegistration, error) {
	dbUserRegistrations, err := bookActivityRegistrationStorage.GetByUserId(userId)

	if err != nil {
		return nil, err
	}

	return dbUserRegistrations.([]*models.BookActivityRegistration), nil
}

func (defaultGameActivityRegistrationService *DefaultGameActivityRegistrationService) GetUserGameActivityRegistrations(userId uint) ([]*models.GameActivityRegistration, error) {
	dbUserRegistrations, err := gameActivityRegistrationStorage.GetByUserId(userId)

	if err != nil {
		return nil, err
	}

	return dbUserRegistrations.([]*models.GameActivityRegistration), nil
}

func (defaultBookActivityRegistrationService *DefaultBookActivityRegistrationService) GetUserBookActivityRegistrationsTimeRange(userId uint, startTime int64, endTime int64) ([]*models.BookActivityRegistration, error) {
	dbUserRegistrations, err := bookActivityRegistrationStorage.GetByUserIdAndTimeRange(userId, startTime, endTime)

	if err != nil {
		return nil, err
	}

	return dbUserRegistrations.([]*models.BookActivityRegistration), nil
}

func (defaultGameActivityRegistrationService *DefaultGameActivityRegistrationService) GetUserGameActivityRegistrationsTimeRange(userId uint, startDate int64, endDate int64) ([]*models.GameActivityRegistration, error) {
	dbUserRegistrations, err := gameActivityRegistrationStorage.GetByUserIdAndInterval(userId, startDate, endDate)

	if err != nil {
		return nil, err
	}

	return dbUserRegistrations.([]*models.GameActivityRegistration), nil
}

func (defaultBookActivityRegistrationService *DefaultBookActivityRegistrationService) CreateBookActivityRegistration(addRegistrationBody *AddBookActivityRegistrationBody) (*models.BookActivityRegistration, error) {
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

func (defaultGameActivityRegistrationService *DefaultGameActivityRegistrationService) CreateGameActivityRegistration(addRegistrationBody *AddGameActivityRegistrationBody) (*models.GameActivityRegistration, error) {

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
