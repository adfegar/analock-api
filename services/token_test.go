package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/adfer-dev/analock-api/models"
	"github.com/stretchr/testify/assert"
)

// mockTokenStorage implements TokenStorageInterface
type mockTokenStorage struct {
	TokensById          map[uint]*models.Token
	TokensByValue       map[string]*models.Token
	TokensByUserAndKind map[string]*models.Token
	TokenPairByUserID   map[uint][2]*models.Token

	GetErr              error
	GetByValueErr       error
	GetByUserAndKindErr error
	GetByUserIdErr      error
	CreateErr           error
	UpdateErr           error
	DeleteErr           error
}

func newMockTokenStorage() *mockTokenStorage {
	return &mockTokenStorage{
		TokensById:          make(map[uint]*models.Token),
		TokensByValue:       make(map[string]*models.Token),
		TokensByUserAndKind: make(map[string]*models.Token),
		TokenPairByUserID:   make(map[uint][2]*models.Token),
	}
}

func (m *mockTokenStorage) Get(id uint) (interface{}, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	token, ok := m.TokensById[id]
	if !ok {
		return nil, errors.New("token not found by id")
	}
	return token, nil
}

func (m *mockTokenStorage) GetByValue(tokenValue string) (interface{}, error) {
	if m.GetByValueErr != nil {
		return nil, m.GetByValueErr
	}
	token, ok := m.TokensByValue[tokenValue]
	if !ok {
		return nil, errors.New("token not found by value")
	}
	return token, nil
}

func (m *mockTokenStorage) GetByUserAndKind(userId uint, kind models.TokenKind) (interface{}, error) {
	if m.GetByUserAndKindErr != nil {
		return nil, m.GetByUserAndKindErr
	}
	key := getTokenStorageKey(userId, kind)
	token, ok := m.TokensByUserAndKind[key]
	if !ok {
		return nil, errors.New("token not found by user and kind")
	}
	return token, nil
}

func (m *mockTokenStorage) GetByUserId(userId uint) ([2]*models.Token, error) {
	if m.GetByUserIdErr != nil {
		return [2]*models.Token{}, m.GetByUserIdErr
	}
	pair, ok := m.TokenPairByUserID[userId]
	if !ok {
		return [2]*models.Token{}, errors.New("token pair not found for user")
	}
	return pair, nil
}

func (m *mockTokenStorage) Create(data interface{}) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	token, ok := data.(*models.Token)
	if !ok {
		return errors.New("create: invalid type for Token")
	}
	if token.Id == 0 {
		token.Id = uint(len(m.TokensById) + 1) // Simple ID generation
	}
	m.TokensById[token.Id] = token
	m.TokensByValue[token.TokenValue] = token
	m.TokensByUserAndKind[getTokenStorageKey(token.UserRefer, token.Kind)] = token
	// For GetByUserId (token pair), we might need more complex logic or assume it's pre-populated for tests needing it.
	return nil
}

func (m *mockTokenStorage) Update(data interface{}) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	token, ok := data.(*models.Token)
	if !ok {
		return errors.New("update: invalid type for Token")
	}
	_, exists := m.TokensById[token.Id]
	if !exists {
		return errors.New("update: token not found")
	}
	m.TokensById[token.Id] = token
	m.TokensByValue[token.TokenValue] = token
	m.TokensByUserAndKind[getTokenStorageKey(token.UserRefer, token.Kind)] = token
	return nil
}

func (m *mockTokenStorage) Delete(id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	token, exists := m.TokensById[id]
	if !exists {
		return errors.New("delete: token not found")
	}
	delete(m.TokensById, id)
	delete(m.TokensByValue, token.TokenValue)
	delete(m.TokensByUserAndKind, getTokenStorageKey(token.UserRefer, token.Kind))
	// Also potentially remove from TokenPairByUserID if managing that explicitly
	return nil
}

// Helper for consistent key generation
func getTokenStorageKey(userId uint, kind models.TokenKind) string {
	return fmt.Sprintf("%d-%d", userId, kind)
}

// --- Test Cases ---
var tokenService TokenService = &DefaultTokenService{}

func TestGetTokenById(t *testing.T) {
	originalTS := tokenStorage
	tokenStorageMock := newMockTokenStorage()
	tokenStorage = tokenStorageMock
	defer func() { tokenStorage = originalTS }()

	testToken := &models.Token{Id: 1, TokenValue: "abc", Kind: models.Access, UserRefer: 10}
	tokenStorageMock.TokensById[testToken.Id] = testToken

	token, err := tokenService.GetTokenById(1)
	assert.NoError(t, err)
	assert.Equal(t, testToken, token)

	_, err = tokenService.GetTokenById(2) // Non-existent
	assert.Error(t, err)

	tokenStorageMock.GetErr = errors.New("forced Get error")
	_, err = tokenService.GetTokenById(1)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Get error")
}

func TestGetTokenByValue(t *testing.T) {
	originalTS := tokenStorage
	tokenStorageMock := newMockTokenStorage()
	tokenStorage = tokenStorageMock
	defer func() { tokenStorage = originalTS }()

	testToken := &models.Token{Id: 1, TokenValue: "token123", Kind: models.Refresh, UserRefer: 11}
	tokenStorageMock.TokensByValue[testToken.TokenValue] = testToken

	token, err := tokenService.GetTokenByValue("token123")
	assert.NoError(t, err)
	assert.Equal(t, testToken, token)

	_, err = tokenService.GetTokenByValue("nonexistent")
	assert.Error(t, err)

	tokenStorageMock.GetByValueErr = errors.New("forced GetByValue error")
	_, err = tokenService.GetTokenByValue("token123")
	assert.Error(t, err)
	assert.EqualError(t, err, "forced GetByValue error")
}

func TestGetUserTokenByKind(t *testing.T) {
	originalTS := tokenStorage
	tokenStorageMock := newMockTokenStorage()
	tokenStorage = tokenStorageMock
	defer func() { tokenStorage = originalTS }()

	userId := uint(15)
	kind := models.Access
	testToken := &models.Token{Id: 5, TokenValue: "kindtoken", Kind: kind, UserRefer: userId}
	tokenStorageMock.TokensByUserAndKind[getTokenStorageKey(userId, kind)] = testToken

	token, err := tokenService.GetUserTokenByKind(userId, kind)
	assert.NoError(t, err)
	assert.Equal(t, testToken, token)

	_, err = tokenService.GetUserTokenByKind(userId, models.Refresh) // Different kind
	assert.Error(t, err)

	tokenStorageMock.GetByUserAndKindErr = errors.New("forced GetByUserAndKind error")
	_, err = tokenService.GetUserTokenByKind(userId, kind)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced GetByUserAndKind error")
}

func TestGetUserTokenPair(t *testing.T) {
	originalTS := tokenStorage
	tokenStorageMock := newMockTokenStorage()
	tokenStorage = tokenStorageMock
	defer func() { tokenStorage = originalTS }()

	userId := uint(20)
	accessToken := &models.Token{Id: 10, TokenValue: "accessPair", Kind: models.Access, UserRefer: userId}
	refreshToken := &models.Token{Id: 11, TokenValue: "refreshPair", Kind: models.Refresh, UserRefer: userId}
	expectedPair := [2]*models.Token{accessToken, refreshToken}
	tokenStorageMock.TokenPairByUserID[userId] = expectedPair

	pair, err := tokenService.GetUserTokenPair(userId)
	assert.NoError(t, err)
	assert.Equal(t, expectedPair, pair)

	_, err = tokenService.GetUserTokenPair(21) // Non-existent user for pair
	assert.Error(t, err)

	tokenStorageMock.GetByUserIdErr = errors.New("forced GetByUserId error")
	_, err = tokenService.GetUserTokenPair(userId)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced GetByUserId error")
}

func TestSaveToken(t *testing.T) {
	originalTS := tokenStorage
	tokenStorageMock := newMockTokenStorage()
	tokenStorage = tokenStorageMock
	defer func() { tokenStorage = originalTS }()

	tokenToSave := &models.Token{TokenValue: "newtoken", Kind: models.Access, UserRefer: 25}

	savedToken, err := tokenService.SaveToken(tokenToSave)
	assert.NoError(t, err)
	assert.NotNil(t, savedToken)
	assert.Equal(t, tokenToSave.TokenValue, savedToken.TokenValue)
	assert.True(t, savedToken.Id > 0) // Check if mock ID was assigned
	assert.Equal(t, tokenToSave, tokenStorageMock.TokensById[savedToken.Id])

	tokenStorageMock.CreateErr = errors.New("forced Create error")
	_, err = tokenService.SaveToken(tokenToSave)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Create error")
}

func TestUpdateToken(t *testing.T) {
	originalTS := tokenStorage
	tokenStorageMock := newMockTokenStorage()
	tokenStorage = tokenStorageMock
	defer func() { tokenStorage = originalTS }()

	initialToken := &models.Token{Id: 30, TokenValue: "initial", Kind: models.Refresh, UserRefer: 30}
	tokenStorageMock.TokensById[initialToken.Id] = initialToken
	tokenStorageMock.TokensByValue[initialToken.TokenValue] = initialToken
	tokenStorageMock.TokensByUserAndKind[getTokenStorageKey(initialToken.UserRefer, initialToken.Kind)] = initialToken

	tokenToUpdate := &models.Token{Id: 30, TokenValue: "updated", Kind: models.Refresh, UserRefer: 30}

	updatedToken, err := tokenService.UpdateToken(tokenToUpdate)
	assert.NoError(t, err)
	assert.Equal(t, tokenToUpdate, updatedToken)
	assert.Equal(t, "updated", tokenStorageMock.TokensById[30].TokenValue)

	tokenStorageMock.UpdateErr = errors.New("forced Update error")
	_, err = tokenService.UpdateToken(tokenToUpdate)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Update error")
}

func TestDeleteToken(t *testing.T) {
	originalTS := tokenStorage
	tokenStorageMock := newMockTokenStorage()
	tokenStorage = tokenStorageMock
	defer func() { tokenStorage = originalTS }()

	tokenToDelete := &models.Token{Id: 40, TokenValue: "deleteme", Kind: models.Access, UserRefer: 40}
	tokenStorageMock.TokensById[tokenToDelete.Id] = tokenToDelete
	tokenStorageMock.TokensByValue[tokenToDelete.TokenValue] = tokenToDelete
	tokenStorageMock.TokensByUserAndKind[getTokenStorageKey(tokenToDelete.UserRefer, tokenToDelete.Kind)] = tokenToDelete

	err := tokenService.DeleteToken(tokenToDelete.Id)
	assert.NoError(t, err)
	_, exists := tokenStorageMock.TokensById[tokenToDelete.Id]
	assert.False(t, exists)

	// Test deleting non-existent
	err = tokenService.DeleteToken(999)
	assert.Error(t, err) // Mock returns error for not found

	tokenStorageMock.DeleteErr = errors.New("forced Delete error")
	// Re-add the token so the delete operation has something to target before the forced error
	tokenStorageMock.TokensById[tokenToDelete.Id] = tokenToDelete
	err = tokenService.DeleteToken(tokenToDelete.Id)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Delete error")
}
