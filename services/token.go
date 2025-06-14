package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

type TokenBody struct {
	TokenValue string           `json:"token" validate:"required,jwt"`
	UserRefer  uint             `json:"user_id" validate:"required,number"`
	Kind       models.TokenKind `json:"kind" validate:"required,number"`
}

var tokenStorage storage.TokenStorageInterface = &storage.TokenStorage{}

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

// TokenServiceImpl is the concrete implementation of TokenService.
type TokenServiceImpl struct{}

// NewTokenServiceImpl creates a new DefaultTokenService.
func NewTokenServiceImpl() *TokenServiceImpl {
	return &TokenServiceImpl{}
}

func (tokenService *TokenServiceImpl) GetTokenById(id uint) (*models.Token, error) {
	token, err := tokenStorage.Get(id)
	if err != nil {
		return nil, err
	}
	return token.(*models.Token), nil
}

func (tokenService *TokenServiceImpl) GetTokenByValue(tokenValue string) (*models.Token, error) {
	token, err := tokenStorage.GetByValue(tokenValue)
	if err != nil {
		return nil, err
	}
	return token.(*models.Token), nil
}

func (tokenService *TokenServiceImpl) GetUserTokenByKind(userId uint, kind models.TokenKind) (*models.Token, error) {
	token, err := tokenStorage.GetByUserAndKind(userId, kind)
	if err != nil {
		return nil, err
	}
	return token.(*models.Token), nil
}

func (tokenService *TokenServiceImpl) GetUserTokenPair(userId uint) ([2]*models.Token, error) {
	tokenPair, err := tokenStorage.GetByUserId(userId)
	if err != nil {
		return [2]*models.Token{}, err
	}
	return tokenPair, nil
}

func (tokenService *TokenServiceImpl) SaveToken(tokenBody *models.Token) (*models.Token, error) {
	err := tokenStorage.Create(tokenBody)
	if err != nil {
		return nil, err
	}
	return tokenBody, nil
}

func (tokenService *TokenServiceImpl) UpdateToken(tokenBody *models.Token) (*models.Token, error) {
	err := tokenStorage.Update(tokenBody)
	if err != nil {
		return nil, err
	}
	return tokenBody, nil
}

func (tokenService *TokenServiceImpl) DeleteToken(id uint) error {
	return tokenStorage.Delete(id)
}
