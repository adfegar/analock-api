package storage

import (
	"database/sql"

	"github.com/adfer-dev/analock-api/database"
	"github.com/adfer-dev/analock-api/models"
)

const (
	getActivityRegistrationByIdentifierQuery = "SELECT * FROM activity_registration WHERE id = ?;"
	insertActivityRegistrationQuery          = "INSERT INTO activity_registration (registration_date, user_id) VALUES (?, ?);"
	updateActivityRegistrationQuery          = "UPDATE activity_registration SET registration_date = ? WHERE id = ?;"
	deleteActivityRegistrationQuery          = "DELETE FROM activity_registration WHERE id = ?;"
)

type ActivityRegistrationStorage struct{}

var activityRegistrationNotFoundError = &models.DbNotFoundError{DbItem: &models.ActivityRegistration{}}
var failedToParseActivityRegistrationError = &models.DbCouldNotParseItemError{DbItem: &models.ActivityRegistration{}}

func (activityRegistrationStorage *ActivityRegistrationStorage) Get(id uint) (interface{}, error) {
	result, err := database.GetDatabaseInstance().GetConnection().Query(getActivityRegistrationByIdentifierQuery, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, activityRegistrationNotFoundError
	}

	scannedDiaryEntry, scanErr := activityRegistrationStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	activityRegistration, ok := scannedDiaryEntry.(models.ActivityRegistration)

	if !ok {
		return nil, failedToParseActivityRegistrationError
	}

	return &activityRegistration, nil
}
func (activityRegistrationStorage *ActivityRegistrationStorage) Create(activityRegistration interface{}) error {
	dbActivityRegistration, ok := activityRegistration.(*models.ActivityRegistration)

	if !ok {
		return failedToParseActivityRegistrationError
	}

	result, err := database.GetDatabaseInstance().GetConnection().Exec(insertActivityRegistrationQuery,
		dbActivityRegistration.RegistrationDate,
		dbActivityRegistration.UserRefer)

	if err != nil {
		return err
	}

	diaryEntryId, idErr := result.LastInsertId()
	if idErr != nil {
		return idErr
	}

	dbActivityRegistration.Id = uint(diaryEntryId)

	return nil
}

func (activityRegistrationStorage *ActivityRegistrationStorage) Update(activityRegistration interface{}) error {
	dbActivityRegistration, ok := activityRegistration.(*models.ActivityRegistration)

	if !ok {
		return failedToParseActivityRegistrationError
	}

	result, err := database.GetDatabaseInstance().GetConnection().Exec(updateActivityRegistrationQuery,
		dbActivityRegistration.RegistrationDate,
		dbActivityRegistration.Id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return activityRegistrationNotFoundError
	}

	return nil
}

func (activityRegistrationStorage *ActivityRegistrationStorage) Delete(id uint) error {

	result, err := database.GetDatabaseInstance().GetConnection().Exec(deleteActivityRegistrationQuery, id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return activityRegistrationNotFoundError
	}

	return nil
}

func (activityRegistrationStorage *ActivityRegistrationStorage) Scan(rows *sql.Rows) (interface{}, error) {
	var activityRegistration models.ActivityRegistration

	scanErr := rows.Scan(&activityRegistration.Id, &activityRegistration.RegistrationDate, &activityRegistration.UserRefer)

	return activityRegistration, scanErr
}
