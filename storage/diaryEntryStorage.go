package storage

import (
	"database/sql"

	"github.com/adfer-dev/analock-api/models"
)

const (
	getDiaryEntryQuery       = "SELECT * FROM diary_entry where id = ?;"
	getUserDiaryEntriesQuery = "SELECT * FROM diary_entry where user_id = ?;"
	insertDiaryEntryQuery    = "INSERT INTO diary_entry (title, content, publishDate, user_id) VALUES (?, ?, ?, ?);"
	updateDiaryEntryQuery    = "UPDATE diary_entry SET title = ?, content = ?, publishDate = ? WHERE id = ?;"
	deleteDiaryEntryQuery    = "DELETE FROM diary_entry WHERE id = ?;"
)

type DiaryEntryStorage struct{}

var diaryEntryNotFoundError = &models.DbNotFoundError{DbItem: &models.DiaryEntry{}}
var failedToParseDiaryEntryError = &models.DbCouldNotParseItemError{DbItem: &models.DiaryEntry{}}

func (diaryEntryStorage *DiaryEntryStorage) Get(id uint) (interface{}, error) {
	result, err := databaseConnection.Query(getDiaryEntryQuery, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, diaryEntryNotFoundError
	}

	scannedDiaryEntry, scanErr := diaryEntryStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	diaryEntry, ok := scannedDiaryEntry.(models.DiaryEntry)

	if !ok {
		return nil, failedToParseDiaryEntryError
	}

	return &diaryEntry, nil
}

func (diaryEntryStorage *DiaryEntryStorage) GetByUserId(userId uint) (interface{}, error) {
	var userDiaryEntries []*models.DiaryEntry
	result, err := databaseConnection.Query(getUserDiaryEntriesQuery, userId)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		scannedDiaryEntry, scanErr := diaryEntryStorage.Scan(result)

		if scanErr != nil {
			return nil, scanErr
		}
		diaryEntry, ok := scannedDiaryEntry.(models.DiaryEntry)

		if !ok {
			return nil, failedToParseDiaryEntryError
		}

		userDiaryEntries = append(userDiaryEntries, &diaryEntry)
	}

	return userDiaryEntries, nil
}

func (diaryEntryStorage *DiaryEntryStorage) Create(diaryEntry interface{}) error {
	dbDiaryEntry, ok := diaryEntry.(*models.DiaryEntry)

	if !ok {
		return failedToParseUserError
	}

	result, err := databaseConnection.Exec(insertDiaryEntryQuery,
		dbDiaryEntry.Title,
		dbDiaryEntry.Content,
		dbDiaryEntry.Date,
		dbDiaryEntry.UserRefer)

	if err != nil {
		return err
	}

	diaryEntryId, idErr := result.LastInsertId()
	if idErr != nil {
		return idErr
	}

	dbDiaryEntry.Id = uint(diaryEntryId)

	return nil
}

func (diaryEntryStorage *DiaryEntryStorage) Update(diaryEntry interface{}) error {
	dbDiaryEntry, ok := diaryEntry.(*models.DiaryEntry)

	if !ok {
		return failedToParseUserError
	}

	result, err := databaseConnection.Exec(updateDiaryEntryQuery,
		dbDiaryEntry.Title,
		dbDiaryEntry.Content,
		dbDiaryEntry.Date,
		dbDiaryEntry.Id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return diaryEntryNotFoundError
	}

	return nil
}

func (diaryEntryStorage *DiaryEntryStorage) Delete(id uint) error {
	result, err := databaseConnection.Exec(deleteDiaryEntryQuery, id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return diaryEntryNotFoundError
	}

	return nil
}

func (diaryEntryStorage *DiaryEntryStorage) Scan(rows *sql.Rows) (interface{}, error) {
	var diaryEntry models.DiaryEntry

	scanErr := rows.Scan(&diaryEntry.Id, &diaryEntry.Title, &diaryEntry.Content, &diaryEntry.Date, &diaryEntry.UserRefer)

	return diaryEntry, scanErr
}
