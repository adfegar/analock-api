package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

var diaryEntryStorage = &storage.DiaryEntryStorage{}

type SaveDiaryEntryBody struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	PublishDate int64  `json:"publishDate" validate:"required"`
	UserRefer   uint   `json:"user_id" validate:"required"`
}

type UpdateDiaryEntryBody struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	PublishDate int64  `json:"publishDate" validate:"required"`
}

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

func SaveDiaryEntry(diaryEntryBody *SaveDiaryEntryBody) (*models.DiaryEntry, error) {
	dbEntry := &models.DiaryEntry{
		Title:     diaryEntryBody.Title,
		Content:   diaryEntryBody.Content,
		Date:      diaryEntryBody.PublishDate,
		UserRefer: diaryEntryBody.UserRefer,
	}
	err := diaryEntryStorage.Create(dbEntry)

	if err != nil {
		return nil, err
	}

	return dbEntry, nil
}

func UpdateDiaryEntry(diaryEntryId uint, diaryEntryBody *UpdateDiaryEntryBody) (*models.DiaryEntry, error) {
	dbEntry := &models.DiaryEntry{
		Id:      diaryEntryId,
		Title:   diaryEntryBody.Title,
		Content: diaryEntryBody.Content,
		Date:    diaryEntryBody.PublishDate,
	}
	err := diaryEntryStorage.Update(dbEntry)

	if err != nil {
		return nil, err
	}

	return dbEntry, nil
}

func DeleteDiaryEntry(id uint) error {
	return diaryEntryStorage.Delete(id)
}
