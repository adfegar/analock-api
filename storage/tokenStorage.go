package storage

import (
	"database/sql"

	"github.com/adfer-dev/analock-api/database"
	"github.com/adfer-dev/analock-api/models"
)

const (
	getTokenQuery              = "SELECT * FROM token where id = ?;"
	getTokenByUserQuery        = "SELECT * FROM token where user_id = ?;"
	getTokenByValueQuery       = "SELECT * FROM token where value = ?;"
	getTokenByUserAndKindQuery = "SELECT * FROM token where user_id = ? AND kind = ?;"
	insertTokenQuery           = "INSERT INTO token (value, kind, user_id) VALUES (?, ?, ?);"
	updateTokenQuery           = "UPDATE token SET value = ?, kind = ? WHERE id = ?;"
	deleteTokenQuery           = "DELETE FROM token WHERE id = ?;"
)

type TokenStorageInterface interface {
	Get(id uint) (interface{}, error)
	GetByValue(tokenValue string) (interface{}, error)
	GetByUserAndKind(userId uint, kind models.TokenKind) (interface{}, error)
	GetByUserId(userId uint) ([2]*models.Token, error)
	Create(data interface{}) error
	Update(data interface{}) error
	Delete(id uint) error
}

type TokenStorage struct{}

var tokenNotFoundError = &models.DbNotFoundError{DbItem: &models.Token{}}
var failedToParseTokenError = &models.DbCouldNotParseItemError{DbItem: &models.Token{}}

func (tokenStorage *TokenStorage) Get(id uint) (interface{}, error) {
	result, err := database.GetDatabaseInstance().GetConnection().Query(getTokenQuery, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, tokenNotFoundError
	}

	scannedToken, scanErr := tokenStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	token, ok := scannedToken.(*models.Token)

	if !ok {
		return nil, failedToParseTokenError
	}

	return token, nil
}

func (tokenStorage *TokenStorage) GetByUserId(id uint) ([2]*models.Token, error) {
	var tokenPair [2]*models.Token
	result, err := database.GetDatabaseInstance().GetConnection().Query(getTokenByUserQuery, id)

	if err != nil {
		return tokenPair, err
	}

	rows := 0
	for result.Next() {
		scannedToken, scanErr := tokenStorage.Scan(result)

		if scanErr != nil {
			return tokenPair, scanErr
		}

		token, ok := scannedToken.(*models.Token)

		if !ok {
			return tokenPair, failedToParseTokenError
		}
		tokenPair[rows] = token
		rows++
	}

	if rows == 0 {
		return tokenPair, tokenNotFoundError
	}

	return tokenPair, nil
}

func (tokenStorage *TokenStorage) GetByValue(tokenValue string) (interface{}, error) {
	result, err := database.GetDatabaseInstance().GetConnection().Query(getTokenByValueQuery, tokenValue)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, tokenNotFoundError
	}

	scannedToken, scanErr := tokenStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	token, ok := scannedToken.(*models.Token)

	if !ok {
		return nil, failedToParseTokenError
	}

	return token, nil
}

func (tokenStorage *TokenStorage) GetByUserAndKind(userId uint, tokenKind models.TokenKind) (interface{}, error) {
	result, err := database.GetDatabaseInstance().GetConnection().Query(getTokenByUserAndKindQuery, userId, tokenKind)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	if !result.Next() {
		return nil, tokenNotFoundError
	}

	scannedToken, scanErr := tokenStorage.Scan(result)

	if scanErr != nil {
		return nil, scanErr
	}

	token, ok := scannedToken.(*models.Token)

	if !ok {
		return nil, failedToParseTokenError
	}

	return token, nil
}

func (tokenStorage *TokenStorage) Create(token interface{}) error {
	dbToken, ok := token.(*models.Token)
	tokenAlreadyExistsError := &models.DbItemAlreadyExistsError{DbItem: &models.Token{}}

	if !ok {
		return failedToParseUserError
	}

	user, getUserErr := tokenStorage.Get(dbToken.Id)
	_, isNotFoundError := getUserErr.(*models.DbNotFoundError)

	if user != nil && !isNotFoundError {
		return tokenAlreadyExistsError
	}

	result, err := database.GetDatabaseInstance().GetConnection().Exec(insertTokenQuery, dbToken.TokenValue, dbToken.Kind, dbToken.UserRefer)
	if err != nil {
		return err
	}

	tokenId, idErr := result.LastInsertId()
	if idErr != nil {
		return idErr
	}

	dbToken.Id = uint(tokenId)

	return nil
}

func (tokenStorage *TokenStorage) Update(token interface{}) error {
	dbToken, ok := token.(*models.Token)

	if !ok {
		return failedToParseUserError
	}

	result, err := database.GetDatabaseInstance().GetConnection().Exec(updateTokenQuery, dbToken.TokenValue, dbToken.Kind, dbToken.Id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return tokenNotFoundError
	}

	return nil
}

func (tokenStorage *TokenStorage) Delete(id uint) error {
	result, err := database.GetDatabaseInstance().GetConnection().Exec(deleteTokenQuery, id)

	if err != nil {
		return err
	}

	affectedRows, errAffectedRows := result.RowsAffected()

	if errAffectedRows != nil {
		return errAffectedRows
	}

	if affectedRows == 0 {
		return tokenNotFoundError
	}

	return nil
}

func (tokenStorage *TokenStorage) Scan(rows *sql.Rows) (interface{}, error) {
	var token models.Token

	scanErr := rows.Scan(&token.Id, &token.TokenValue, &token.Kind, &token.UserRefer)

	return &token, scanErr
}
