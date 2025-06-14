package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

// Interface for DiaryEntryStorage
type DiaryEntryStorageInterface interface {
	Get(id uint) (interface{}, error)
	GetByUserId(userId uint) (interface{}, error)
	GetByUserIdAndDateInterval(userId uint, startDate int64, endDate int64) (interface{}, error)
	Create(data interface{}) error
	Update(data interface{}) error
}

// ActivityRegistrationStorageInterface is defined in activityRegistration.go
// and is accessible throughout the 'services' package.

type SaveDiaryEntryBody struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	PublishDate int64  `json:"publishDate" validate:"required"`
	UserRefer   uint   `json:"userId" validate:"required"`
}

// Added back UpdateDiaryEntryBody
type UpdateDiaryEntryBody struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	PublishDate int64  `json:"publishDate" validate:"required"`
}

var diaryEntryStorage DiaryEntryStorageInterface = &storage.DiaryEntryStorage{}

type DiaryEntryService interface {
	GetDiaryEntryById(id uint) (*models.DiaryEntry, error)
	GetUserEntries(userId uint) ([]*models.DiaryEntry, error)
	GetUserEntriesTimeRange(userId uint, startDate int64, endDate int64) ([]*models.DiaryEntry, error)
	SaveDiaryEntry(diaryEntryBody *SaveDiaryEntryBody) (*models.DiaryEntry, error)
	UpdateDiaryEntry(diaryEntryId uint, diaryEntryBody *UpdateDiaryEntryBody) (*models.DiaryEntry, error)
	DeleteDiaryEntry(id uint) error
}

type DefaultDiaryEntryService struct{}

func (defaultDiaryEntryService *DefaultDiaryEntryService) GetDiaryEntryById(id uint) (*models.DiaryEntry, error) {
	diaryEntry, err := diaryEntryStorage.Get(id)

	if err != nil {
		return nil, err
	}

	return diaryEntry.(*models.DiaryEntry), nil
}

func (defaultDiaryEntryService *DefaultDiaryEntryService) GetUserEntries(userId uint) ([]*models.DiaryEntry, error) {

	diaryEntry, err := diaryEntryStorage.GetByUserId(userId)

	if err != nil {
		return nil, err
	}
	return diaryEntry.([]*models.DiaryEntry), nil
}

func (defaultDiaryEntryService *DefaultDiaryEntryService) GetUserEntriesTimeRange(userId uint, startDate int64, endDate int64) ([]*models.DiaryEntry, error) {
	diaryEntry, err := diaryEntryStorage.GetByUserIdAndDateInterval(userId, startDate, endDate)

	if err != nil {
		return nil, err
	}

	return diaryEntry.([]*models.DiaryEntry), nil
}

func (defaultDiaryEntryService *DefaultDiaryEntryService) SaveDiaryEntry(diaryEntryBody *SaveDiaryEntryBody) (*models.DiaryEntry, error) {
	dbActivityRegistration := &models.ActivityRegistration{
		RegistrationDate: diaryEntryBody.PublishDate,
		UserRefer:        diaryEntryBody.UserRefer,
	}

	saveRegistrationErr := activityRegistrationStorage.Create(dbActivityRegistration)

	if saveRegistrationErr != nil {
		return nil, saveRegistrationErr
	}

	dbEntry := &models.DiaryEntry{
		Title:        diaryEntryBody.Title,
		Content:      diaryEntryBody.Content,
		Registration: *dbActivityRegistration,
	}
	err := diaryEntryStorage.Create(dbEntry)

	if err != nil {
		return nil, err
	}

	return dbEntry, nil
}

func (defaultDiaryEntryService *DefaultDiaryEntryService) UpdateDiaryEntry(diaryEntryId uint, diaryEntryBody *UpdateDiaryEntryBody) (*models.DiaryEntry, error) {
	storedDiaryEntry, getDiaryEntryError := defaultDiaryEntryService.GetDiaryEntryById(diaryEntryId)

	if getDiaryEntryError != nil {
		return nil, getDiaryEntryError
	}

	dbRegistration := &models.ActivityRegistration{
		Id:               storedDiaryEntry.Registration.Id,
		RegistrationDate: diaryEntryBody.PublishDate,
		UserRefer:        storedDiaryEntry.Registration.UserRefer,
	}
	updateRegistrationErr := activityRegistrationStorage.Update(dbRegistration)

	if updateRegistrationErr != nil {
		return nil, updateRegistrationErr
	}

	updatedDiaryEntry := &models.DiaryEntry{
		Id:           diaryEntryId,
		Title:        diaryEntryBody.Title,
		Content:      diaryEntryBody.Content,
		Registration: *dbRegistration,
	}
	err := diaryEntryStorage.Update(updatedDiaryEntry)

	if err != nil {
		return nil, err
	}

	return updatedDiaryEntry, nil
}

func (defaultDiaryEntryService *DefaultDiaryEntryService) DeleteDiaryEntry(id uint) error {
	diaryEntry, err := defaultDiaryEntryService.GetDiaryEntryById(id)

	if err != nil {
		return err
	}

	return activityRegistrationStorage.Delete(diaryEntry.Registration.Id)
}
