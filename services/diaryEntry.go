package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

type SaveDiaryEntryBody struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	PublishDate int64  `json:"publishDate" validate:"required"`
	UserRefer   uint   `json:"userId" validate:"required"`
}

type UpdateDiaryEntryBody struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	PublishDate int64  `json:"publishDate" validate:"required"`
}

var diaryEntryStorage = &storage.DiaryEntryStorage{}
var activityRegistrationStorage = &storage.ActivityRegistrationStorage{}

func GetDiaryEntryById(id uint) (*models.DiaryEntry, error) {
	diaryEntry, err := diaryEntryStorage.Get(id)

	if err != nil {
		return nil, err
	}

	return diaryEntry.(*models.DiaryEntry), nil
}

func GetUserEntries(userId uint) ([]*models.DiaryEntry, error) {
	diaryEntry, err := diaryEntryStorage.GetByUserId(userId)

	if err != nil {
		return nil, err
	}

	return diaryEntry.([]*models.DiaryEntry), nil
}

func GetUserEntriesTimeRange(userId uint, startDate int64, endDate int64) ([]*models.DiaryEntry, error) {
	diaryEntry, err := diaryEntryStorage.GetByUserIdAndDateInterval(userId, startDate, endDate)

	if err != nil {
		return nil, err
	}

	return diaryEntry.([]*models.DiaryEntry), nil
}

func SaveDiaryEntry(diaryEntryBody *SaveDiaryEntryBody) (*models.DiaryEntry, error) {
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

func UpdateDiaryEntry(diaryEntryId uint, diaryEntryBody *UpdateDiaryEntryBody) (*models.DiaryEntry, error) {

	storedDiaryEntry, getDiaryEntryError := GetDiaryEntryById(diaryEntryId)

	if getDiaryEntryError != nil {
		return nil, getDiaryEntryError
	}

	dbRegistration := &models.ActivityRegistration{
		Id:               storedDiaryEntry.Registration.Id,
		RegistrationDate: diaryEntryBody.PublishDate,
	}
	updateRegistrationErr := activityRegistrationStorage.Update(dbRegistration)

	if updateRegistrationErr != nil {
		return nil, updateRegistrationErr
	}

	updatedDiaryEntry := &models.DiaryEntry{
		Id:      diaryEntryId,
		Title:   diaryEntryBody.Title,
		Content: diaryEntryBody.Content,
		Registration: models.ActivityRegistration{
			Id:               dbRegistration.Id,
			RegistrationDate: dbRegistration.RegistrationDate,
			UserRefer:        storedDiaryEntry.Registration.UserRefer,
		},
	}
	err := diaryEntryStorage.Update(updatedDiaryEntry)

	if err != nil {
		return nil, err
	}

	return updatedDiaryEntry, nil
}

func DeleteDiaryEntry(id uint) error {
	diaryEntry, err := GetDiaryEntryById(id)

	if err != nil {
		return err
	}

	return activityRegistrationStorage.Delete(diaryEntry.Registration.Id)
}
