package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

type TokenBody struct {
	TokenValue string           `json:"token" validate:"required,jwt"`
	UserRefer  uint             `json:"user_id" validate:"required,number"`
	Kind       models.TokenKind `validate:"required,number"`
}

var tokenStorage *storage.TokenStorage = &storage.TokenStorage{}

func GetTokenById(id uint) (*models.Token, error) {
	token, err := tokenStorage.Get(id)

	return token.(*models.Token), err
}

func GetTokenByValue(tokenValue string) (*models.Token, error) {
	token, err := tokenStorage.GetByValue(tokenValue)

	return token.(*models.Token), err
}

func GetUserTokenByKind(userId uint, kind models.TokenKind) (*models.Token, error) {
	token, err := tokenStorage.GetByUserAndKind(userId, kind)

	return token.(*models.Token), err
}

func GetUserTokenPair(userId uint) ([2]*models.Token, error) {
	tokenPair, err := tokenStorage.GetByUserId(userId)

	return tokenPair, err
}

func SaveToken(tokenBody *models.Token) (*models.Token, error) {
	err := tokenStorage.Create(tokenBody)

	if err != nil {
		return nil, err
	}

	return tokenBody, nil
}

func UpdateToken(tokenBody *models.Token) (*models.Token, error) {
	err := tokenStorage.Update(tokenBody)

	if err != nil {
		return nil, err
	}

	return tokenBody, nil
}

func DeleteToken(id uint) error {
	return tokenStorage.Delete(id)
}
