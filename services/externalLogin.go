package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

// ExternalLoginStorageInterface defines storage operations for external logins.
type ExternalLoginStorageInterface interface {
	Get(id uint) (interface{}, error)
	GetByClientId(clientId string) (interface{}, error)
	Create(data interface{}) error
	Update(data interface{}) error
	UpdateUserExternalLoginToken(data interface{}) error
	Delete(id uint) error
}

var externalLoginStorage ExternalLoginStorageInterface = &storage.ExternalLoginStorage{}

// ExternalLoginService defines all operations for the external login service.
type ExternalLoginService interface {
	GetExternalLoginById(id uint) (*models.ExternalLogin, error)
	GetExternalLoginByClientId(clientId string) (*models.ExternalLogin, error)
	SaveExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error)
	UpdateExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error)
	UpdateUserExternalLoginToken(userId uint, externalLoginBody *UpdateExternalLoginBody) (*models.ExternalLogin, error)
	DeleteExternalLogin(id uint) error
}

// DefaultExternalLoginService is the concrete implementation of ExternalLoginService.
type DefaultExternalLoginService struct{}

// NewDefaultExternalLoginService creates a new DefaultExternalLoginService.
func NewDefaultExternalLoginService() *DefaultExternalLoginService {
	return &DefaultExternalLoginService{}
}

func (s *DefaultExternalLoginService) GetExternalLoginById(id uint) (*models.ExternalLogin, error) {
	externalLogin, err := externalLoginStorage.Get(id)
	if err != nil {
		return nil, err
	}
	return externalLogin.(*models.ExternalLogin), nil
}

func (s *DefaultExternalLoginService) GetExternalLoginByClientId(clientId string) (*models.ExternalLogin, error) {
	externalLogin, err := externalLoginStorage.GetByClientId(clientId)
	if err != nil {
		return nil, err
	}
	return externalLogin.(*models.ExternalLogin), nil
}

func (s *DefaultExternalLoginService) SaveExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error) {
	err := externalLoginStorage.Create(externalLoginBody)
	if err != nil {
		return nil, err
	}
	return externalLoginBody, nil
}

func (s *DefaultExternalLoginService) UpdateExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error) {
	err := externalLoginStorage.Update(externalLoginBody)
	if err != nil {
		return nil, err
	}
	return externalLoginBody, nil
}

func (s *DefaultExternalLoginService) UpdateUserExternalLoginToken(userId uint, externalLoginBody *UpdateExternalLoginBody) (*models.ExternalLogin, error) {
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

func (s *DefaultExternalLoginService) DeleteExternalLogin(id uint) error {
	return externalLoginStorage.Delete(id)
}

// --- Original Package-Level Functions (Logic moved to DefaultExternalLoginService methods) ---

// func GetExternalLoginById(id uint) (*models.ExternalLogin, error) { ... }
// func GetExternalLoginByClientId(clientId string) (*models.ExternalLogin, error) { ... }
// func SaveExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error) { ... }
// func UpdateExternalLogin(externalLoginBody *models.ExternalLogin) (*models.ExternalLogin, error) { ... }
// func UpdateUserExternalLoginToken(userId uint, externalLoginBody *UpdateExternalLoginBody) (*models.ExternalLogin, error) { ... }
// func DeleteExternalLogin(id uint) error { ... }

// --- Request/Response Bodies ---
type UpdateExternalLoginBody struct {
	ClientToken string `json:"provider_client_token"`
}
