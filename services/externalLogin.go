package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

type UpdateExternalLoginBody struct {
	ClientToken string `json:"provider_client_token"`
}

var externalLoginStorage storage.ExternalLoginStorageInterface = &storage.ExternalLoginStorage{}

// ExternalLoginService defines all operations for the external login service.
type ExternalLoginService interface {
	GetExternalLoginById(id uint) (*models.ExternalLogin, error)
	GetExternalLoginByClientId(clientId string) (*models.ExternalLogin, error)
	SaveExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error)
	UpdateExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error)
	UpdateUserExternalLoginToken(userId uint, externalLoginBody *UpdateExternalLoginBody) (*models.ExternalLogin, error)
	DeleteExternalLogin(id uint) error
}

// ExternalLoginServiceImpl is the concrete implementation of ExternalLoginService.
type ExternalLoginServiceImpl struct{}

// NewExternalLoginServiceImpl creates a new DefaultExternalLoginService.
func NewExternalLoginServiceImpl() *ExternalLoginServiceImpl {
	return &ExternalLoginServiceImpl{}
}

func (externalLoginService *ExternalLoginServiceImpl) GetExternalLoginById(id uint) (*models.ExternalLogin, error) {
	externalLogin, err := externalLoginStorage.Get(id)
	if err != nil {
		return nil, err
	}
	return externalLogin.(*models.ExternalLogin), nil
}

func (externalLoginService *ExternalLoginServiceImpl) GetExternalLoginByClientId(clientId string) (*models.ExternalLogin, error) {
	externalLogin, err := externalLoginStorage.GetByClientId(clientId)
	if err != nil {
		return nil, err
	}
	return externalLogin.(*models.ExternalLogin), nil
}

func (externalLoginService *ExternalLoginServiceImpl) SaveExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error) {
	err := externalLoginStorage.Create(externalLoginBody)
	if err != nil {
		return nil, err
	}
	return externalLoginBody, nil
}

func (externalLoginService *ExternalLoginServiceImpl) UpdateExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error) {
	err := externalLoginStorage.Update(externalLoginBody)
	if err != nil {
		return nil, err
	}
	return externalLoginBody, nil
}

func (externalLoginService *ExternalLoginServiceImpl) UpdateUserExternalLoginToken(userId uint, externalLoginBody *UpdateExternalLoginBody) (*models.ExternalLogin, error) {
	dbExternalLogin := &models.ExternalLogin{
		UserRefer:   userId,
		ClientToken: externalLoginBody.ClientToken,
	}
	err := externalLoginStorage.UpdateUserExternalLoginToken(dbExternalLogin)
	if err != nil {
		return nil, err
	}
	return dbExternalLogin, nil
}

func (externalLoginService *ExternalLoginServiceImpl) DeleteExternalLogin(id uint) error {
	return externalLoginStorage.Delete(id)
}
