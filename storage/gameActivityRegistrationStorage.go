package storage

import (
	"database/sql"

	"github.com/adfer-dev/analock-api/models"
)

const (
	getGameActivityRegistrationByIdentifierQuery = "SELECT arg.id, arb.game_name, ar.registration_date, ar.user_id FROM activity_registration_game arg INNER JOIN activity_registration ar ON (arg.registration_id = ar.id) WHERE arg.id = ?;"
	getUserGameActivityRegistrationsQuery        = "SELECT arg.id, arg.game_name, ar.registration_date, ar.user_id FROM activity_registration_game arg INNER JOIN activity_registration ar ON (arg.registration_id = ar.id) WHERE ar.user_id = ?;"
	insertGameActivityRegistrationQuery          = "INSERT INTO activity_registration_game (game_name, registration_id) VALUES (?, ?);"
	updateGameActivityRegistrationQuery          = "UPDATE activity_registration_game SET game_name = ? WHERE id = ?;"
	deleteGameActivityRegistrationQuery          = "DELETE FROM activity_registration_game WHERE id = ?;"
)

type GameActivityRegistrationStorage struct{}

var gameActivityRegistrationNotFoundError = &models.DbNotFoundError{DbItem: &models.GameActivityRegistration{}}
var failedToParseGameActivityRegistrationError = &models.DbCouldNotParseItemError{DbItem: &models.GameActivityRegistration{}}

func (gameActivityRegistrationStorage *GameActivityRegistrationStorage) Get(id uint) (interface{}, error) {
	result, err := databaseConnection.Query(getGameActivityRegistrationByIdentifierQuery, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, gameActivityRegistrationNotFoundError
	}

	scannedGameActivityRegistration, scanErr := gameActivityRegistrationStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	gameActivityRegistration, ok := scannedGameActivityRegistration.(models.GameActivityRegistration)

	if !ok {
		return nil, failedToParseGameActivityRegistrationError
	}

	return &gameActivityRegistration, nil
}

func (gameActivityRegistrationStorage *GameActivityRegistrationStorage) GetByUserId(userId uint) (interface{}, error) {
	var userGameActivityRegistrations []*models.GameActivityRegistration
	result, err := databaseConnection.Query(getUserGameActivityRegistrationsQuery, userId)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		scannedGameActivityRegistration, scanErr := gameActivityRegistrationStorage.Scan(result)

		if scanErr != nil {
			return nil, scanErr
		}
		gameActivityRegistration, ok := scannedGameActivityRegistration.(models.GameActivityRegistration)

		if !ok {
			return nil, failedToParseGameActivityRegistrationError
		}

		userGameActivityRegistrations = append(userGameActivityRegistrations, &gameActivityRegistration)
	}

	return userGameActivityRegistrations, nil
}

func (gameActivityRegistrationStorage *GameActivityRegistrationStorage) Create(gameRegistration interface{}) error {
	dbGameRegistration, ok := gameRegistration.(*models.GameActivityRegistration)

	if !ok {
		return failedToParseDiaryEntryError
	}

	result, err := databaseConnection.Exec(insertGameActivityRegistrationQuery,
		dbGameRegistration.GameName,
		dbGameRegistration.Registration.Id)

	if err != nil {
		return err
	}

	gameRegistrationId, idErr := result.LastInsertId()
	if idErr != nil {
		return idErr
	}

	dbGameRegistration.Id = uint(gameRegistrationId)

	return nil
}

func (gameActivityRegistrationStorage *GameActivityRegistrationStorage) Update(gameRegistration interface{}) error {
	dbGameRegistration, ok := gameRegistration.(*models.GameActivityRegistration)

	if !ok {
		return failedToParseGameActivityRegistrationError
	}

	result, err := databaseConnection.Exec(updateDiaryEntryQuery,
		dbGameRegistration.GameName,
		dbGameRegistration.Id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return gameActivityRegistrationNotFoundError
	}

	return nil
}

func (gameActivityRegistrationStorage *GameActivityRegistrationStorage) Delete(id uint) error {
	result, err := databaseConnection.Exec(deleteGameActivityRegistrationQuery, id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return gameActivityRegistrationNotFoundError
	}

	return nil
}

func (gameActivityRegistrationStorage *GameActivityRegistrationStorage) Scan(rows *sql.Rows) (interface{}, error) {
	var gameActivityRegistration models.GameActivityRegistration

	scanErr := rows.Scan(&gameActivityRegistration.Id, &gameActivityRegistration.GameName,
		&gameActivityRegistration.Registration.RegistrationDate, &gameActivityRegistration.Registration.UserRefer)

	return gameActivityRegistration, scanErr
}
