package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

// TokenStorageInterface defines storage operations for tokens.
type TokenStorageInterface interface {
	Get(id uint) (interface{}, error)
	GetByValue(tokenValue string) (interface{}, error)
	GetByUserAndKind(userId uint, kind models.TokenKind) (interface{}, error)
	GetByUserId(userId uint) ([2]*models.Token, error)
	Create(data interface{}) error
	Update(data interface{}) error
	Delete(id uint) error
}

var tokenStorage TokenStorageInterface = &storage.TokenStorage{}

// TokenService defines all operations for the token service.
type TokenService interface {
	GetTokenById(id uint) (*models.Token, error)
	GetTokenByValue(tokenValue string) (*models.Token, error)
	GetUserTokenByKind(userId uint, kind models.TokenKind) (*models.Token, error)
	GetUserTokenPair(userId uint) ([2]*models.Token, error)
	SaveToken(tokenBody *models.Token) (*models.Token, error)
	UpdateToken(tokenBody *models.Token) (*models.Token, error)
	DeleteToken(id uint) error
}

// DefaultTokenService is the concrete implementation of TokenService.
type DefaultTokenService struct{}

// NewDefaultTokenService creates a new DefaultTokenService.
func NewDefaultTokenService() *DefaultTokenService {
	return &DefaultTokenService{}
}

func (s *DefaultTokenService) GetTokenById(id uint) (*models.Token, error) {
	token, err := tokenStorage.Get(id)
	if err != nil {
		return nil, err
	}
	return token.(*models.Token), nil
}

func (s *DefaultTokenService) GetTokenByValue(tokenValue string) (*models.Token, error) {
	token, err := tokenStorage.GetByValue(tokenValue)
	if err != nil {
		return nil, err
	}
	return token.(*models.Token), nil
}

func (s *DefaultTokenService) GetUserTokenByKind(userId uint, kind models.TokenKind) (*models.Token, error) {
	token, err := tokenStorage.GetByUserAndKind(userId, kind)
	if err != nil {
		return nil, err
	}
	return token.(*models.Token), nil
}

func (s *DefaultTokenService) GetUserTokenPair(userId uint) ([2]*models.Token, error) {
	tokenPair, err := tokenStorage.GetByUserId(userId)
	if err != nil {
		return [2]*models.Token{}, err
	}
	return tokenPair, nil
}

func (s *DefaultTokenService) SaveToken(tokenBody *models.Token) (*models.Token, error) {
	err := tokenStorage.Create(tokenBody)
	if err != nil {
		return nil, err
	}
	return tokenBody, nil
}

func (s *DefaultTokenService) UpdateToken(tokenBody *models.Token) (*models.Token, error) {
	err := tokenStorage.Update(tokenBody)
	if err != nil {
		return nil, err
	}
	return tokenBody, nil
}

func (s *DefaultTokenService) DeleteToken(id uint) error {
	return tokenStorage.Delete(id)
}

// --- Original Package-Level Functions (Logic moved to DefaultTokenService methods) ---

// func GetTokenById(id uint) (*models.Token, error) { ... }
// func GetTokenByValue(tokenValue string) (*models.Token, error) { ... }
// func GetUserTokenByKind(userId uint, kind models.TokenKind) (*models.Token, error) { ... }
// func GetUserTokenPair(userId uint) ([2]*models.Token, error) { ... }
// func SaveToken(tokenBody *models.Token) (*models.Token, error) { ... }
// func UpdateToken(tokenBody *models.Token) (*models.Token, error) { ... }
// func DeleteToken(id uint) error { ... }

// --- Request/Response Bodies ---

type TokenBody struct {
	TokenValue string           `json:"token" validate:"required,jwt"`
	UserRefer  uint             `json:"user_id" validate:"required,number"`
	Kind       models.TokenKind `json:"kind" validate:"required,number"`
}
