package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

var externalLoginStorage *storage.ExternalLoginStorage = &storage.ExternalLoginStorage{}

func GetExternalLoginById(id uint) (*models.ExternalLogin, error) {
	externalLogin, err := externalLoginStorage.Get(id)

	if err != nil {
		return nil, err
	}

	return externalLogin.(*models.ExternalLogin), nil
}

func GetExternalLoginByClientId(clientId string) (*models.ExternalLogin, error) {
	externalLogin, err := externalLoginStorage.GetByClientId(clientId)

	if err != nil {
		return nil, err
	}

	return externalLogin.(*models.ExternalLogin), nil
}

func SaveExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error) {
	err := externalLoginStorage.Create(externalLoginBody)

	if err != nil {
		return nil, err
	}

	return externalLoginBody, nil
}

func UpdateExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error) {
	err := externalLoginStorage.Update(externalLoginBody)

	if err != nil {
		return nil, err
	}

	return externalLoginBody, nil
}

func DeleteExternalLogin(id uint) error {
	return externalLoginStorage.Delete(id)
}
