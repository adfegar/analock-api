package services

import (
	"errors"
	"testing"

	"github.com/adfer-dev/analock-api/models"
	"github.com/stretchr/testify/assert"
)

// mockExternalLoginStorage implements ExternalLoginStorageInterface
type mockExternalLoginStorage struct {
	LoginsById                      map[uint]*models.ExternalLogin
	LoginsByClientId                map[string]*models.ExternalLogin
	GetErr                          error
	GetByClientIdErr                error
	CreateErr                       error
	UpdateErr                       error
	UpdateUserExternalLoginTokenErr error
	DeleteErr                       error
	LastUpdatedUserTokenLogin       *models.ExternalLogin
}

func newMockExternalLoginStorage() *mockExternalLoginStorage {
	return &mockExternalLoginStorage{
		LoginsById:       make(map[uint]*models.ExternalLogin),
		LoginsByClientId: make(map[string]*models.ExternalLogin),
	}
}

func (externalLoginStorageMock *mockExternalLoginStorage) Get(id uint) (interface{}, error) {
	if externalLoginStorageMock.GetErr != nil {
		return nil, externalLoginStorageMock.GetErr
	}
	login, ok := externalLoginStorageMock.LoginsById[id]
	if !ok {
		return nil, errors.New("external login not found by id")
	}
	return login, nil
}

func (externalLoginStorageMock *mockExternalLoginStorage) GetByClientId(clientId string) (interface{}, error) {
	if externalLoginStorageMock.GetByClientIdErr != nil {
		return nil, externalLoginStorageMock.GetByClientIdErr
	}
	login, ok := externalLoginStorageMock.LoginsByClientId[clientId]
	if !ok {
		return nil, errors.New("external login not found by client id")
	}
	return login, nil
}

func (externalLoginStorageMock *mockExternalLoginStorage) Create(data interface{}) error {
	if externalLoginStorageMock.CreateErr != nil {
		return externalLoginStorageMock.CreateErr
	}
	login, ok := data.(*models.ExternalLogin)
	if !ok {
		return errors.New("create: invalid type for ExternalLogin")
	}
	if login.Id == 0 {
		login.Id = uint(len(externalLoginStorageMock.LoginsById) + 1)
	}
	externalLoginStorageMock.LoginsById[login.Id] = login
	externalLoginStorageMock.LoginsByClientId[login.ClientId] = login
	return nil
}

func (externalLoginStorageMock *mockExternalLoginStorage) Update(data interface{}) error {
	if externalLoginStorageMock.UpdateErr != nil {
		return externalLoginStorageMock.UpdateErr
	}
	login, ok := data.(*models.ExternalLogin)
	if !ok {
		return errors.New("update: invalid type for ExternalLogin")
	}
	_, exists := externalLoginStorageMock.LoginsById[login.Id]
	if !exists {
		return errors.New("update: external login not found")
	}
	externalLoginStorageMock.LoginsById[login.Id] = login
	externalLoginStorageMock.LoginsByClientId[login.ClientId] = login
	return nil
}

func (externalLoginStorageMock *mockExternalLoginStorage) UpdateUserExternalLoginToken(data interface{}) error {
	if externalLoginStorageMock.UpdateUserExternalLoginTokenErr != nil {
		return externalLoginStorageMock.UpdateUserExternalLoginTokenErr
	}
	login, ok := data.(*models.ExternalLogin)
	if !ok {
		return errors.New("updateuserexternallogintoken: invalid type for ExternalLogin")
	}
	// For testing, we can just store what was passed or simulate an update
	// Here, we'll store it to check the UserRefer and ClientToken passed from the service.
	externalLoginStorageMock.LastUpdatedUserTokenLogin = login
	// A more complete mock might try to find an existing login by UserRefer and update its token.
	return nil
}

func (externalLoginStorageMock *mockExternalLoginStorage) Delete(id uint) error {
	if externalLoginStorageMock.DeleteErr != nil {
		return externalLoginStorageMock.DeleteErr
	}
	login, exists := externalLoginStorageMock.LoginsById[id]
	if !exists {
		return errors.New("delete: external login not found")
	}
	delete(externalLoginStorageMock.LoginsById, id)
	delete(externalLoginStorageMock.LoginsByClientId, login.ClientId)
	return nil
}

// --- Test Cases ---
var externalLoginService ExternalLoginService = &ExternalLoginServiceImpl{}

func TestGetExternalLoginById(t *testing.T) {
	originalExternalLoginStorage := externalLoginStorage
	externalLoginStorageMock := newMockExternalLoginStorage()
	externalLoginStorage = externalLoginStorageMock
	defer func() { externalLoginStorage = originalExternalLoginStorage }()

	testLogin := &models.ExternalLogin{Id: 1, ClientId: "client1", UserRefer: 10}
	externalLoginStorageMock.LoginsById[testLogin.Id] = testLogin

	login, err := externalLoginService.GetExternalLoginById(1)
	assert.NoError(t, err)
	assert.Equal(t, testLogin, login)

	_, err = externalLoginService.GetExternalLoginById(2) // Non-existent
	assert.Error(t, err)

	externalLoginStorageMock.GetErr = errors.New("forced Get error")
	_, err = externalLoginService.GetExternalLoginById(1)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Get error")
}

func TestGetExternalLoginByClientId(t *testing.T) {
	originalExternalLoginStorage := externalLoginStorage
	externalLoginStorageMock := newMockExternalLoginStorage()
	externalLoginStorage = externalLoginStorageMock
	defer func() { externalLoginStorage = originalExternalLoginStorage }()

	testLogin := &models.ExternalLogin{Id: 1, ClientId: "client-abc", UserRefer: 11}
	externalLoginStorageMock.LoginsByClientId[testLogin.ClientId] = testLogin

	login, err := externalLoginService.GetExternalLoginByClientId("client-abc")
	assert.NoError(t, err)
	assert.Equal(t, testLogin, login)

	_, err = externalLoginService.GetExternalLoginByClientId("nonexistent-client")
	assert.Error(t, err)

	externalLoginStorageMock.GetByClientIdErr = errors.New("forced GetByClientId error")
	_, err = externalLoginService.GetExternalLoginByClientId("client-abc")
	assert.Error(t, err)
	assert.EqualError(t, err, "forced GetByClientId error")
}

func TestSaveExternalLogin(t *testing.T) {
	originalExternalLoginStorage := externalLoginStorage
	externalLoginStorageMock := newMockExternalLoginStorage()
	externalLoginStorage = externalLoginStorageMock
	defer func() { externalLoginStorage = originalExternalLoginStorage }()

	loginToSave := &models.ExternalLogin{ClientId: "new-client", UserRefer: 12, ClientToken: "token"}

	savedLogin, err := externalLoginService.SaveExternalLogin(loginToSave)
	assert.NoError(t, err)
	assert.NotNil(t, savedLogin)
	assert.True(t, savedLogin.Id > 0) // Mock should assign an ID
	assert.Equal(t, loginToSave.ClientId, savedLogin.ClientId)
	assert.Equal(t, loginToSave, externalLoginStorageMock.LoginsById[savedLogin.Id])

	externalLoginStorageMock.CreateErr = errors.New("forced Create error")
	_, err = externalLoginService.SaveExternalLogin(loginToSave)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Create error")
}

func TestUpdateExternalLogin(t *testing.T) {
	originalExternalLoginStorage := externalLoginStorage
	externalLoginStorageMock := newMockExternalLoginStorage()
	externalLoginStorage = externalLoginStorageMock
	defer func() { externalLoginStorage = originalExternalLoginStorage }()

	initialLogin := &models.ExternalLogin{Id: 20, ClientId: "client-initial", UserRefer: 15}
	externalLoginStorageMock.LoginsById[initialLogin.Id] = initialLogin
	externalLoginStorageMock.LoginsByClientId[initialLogin.ClientId] = initialLogin

	loginToUpdate := &models.ExternalLogin{Id: 20, ClientId: "client-updated", UserRefer: 15, ClientToken: "new-token"}

	updatedLogin, err := externalLoginService.UpdateExternalLogin(loginToUpdate)
	assert.NoError(t, err)
	assert.Equal(t, loginToUpdate, updatedLogin)
	assert.Equal(t, "client-updated", externalLoginStorageMock.LoginsById[20].ClientId)

	externalLoginStorageMock.UpdateErr = errors.New("forced Update error")
	_, err = externalLoginService.UpdateExternalLogin(loginToUpdate)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Update error")
}

func TestUpdateUserExternalLoginToken(t *testing.T) {
	originalExternalLoginStorage := externalLoginStorage
	externalLoginStorageMock := newMockExternalLoginStorage()
	externalLoginStorage = externalLoginStorageMock
	defer func() { externalLoginStorage = originalExternalLoginStorage }()

	userId := uint(25)
	updateBody := &UpdateExternalLoginBody{ClientToken: "new-user-token"}

	updatedLogin, err := externalLoginService.UpdateUserExternalLoginToken(userId, updateBody)

	assert.NoError(t, err)
	assert.NotNil(t, updatedLogin)
	// The service function constructs a new ExternalLogin and passes it to storage.
	// So, the returned updatedLogin is what was constructed by the service.
	assert.Equal(t, userId, updatedLogin.UserRefer)
	assert.Equal(t, updateBody.ClientToken, updatedLogin.ClientToken)

	// Check what was passed to the mock storage's UpdateUserExternalLoginToken method
	assert.NotNil(t, externalLoginStorageMock.LastUpdatedUserTokenLogin)
	assert.Equal(t, userId, externalLoginStorageMock.LastUpdatedUserTokenLogin.UserRefer)
	assert.Equal(t, updateBody.ClientToken, externalLoginStorageMock.LastUpdatedUserTokenLogin.ClientToken)

	externalLoginStorageMock.UpdateUserExternalLoginTokenErr = errors.New("forced UpdateUserExternalLoginToken error")
	_, err = externalLoginService.UpdateUserExternalLoginToken(userId, updateBody)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced UpdateUserExternalLoginToken error")
}

func TestDeleteExternalLogin(t *testing.T) {
	originalELS := externalLoginStorage
	externalLoginStorageMock := newMockExternalLoginStorage()
	externalLoginStorage = externalLoginStorageMock
	defer func() { externalLoginStorage = originalELS }()

	loginToDelete := &models.ExternalLogin{Id: 30, ClientId: "client-delete"}
	externalLoginStorageMock.LoginsById[loginToDelete.Id] = loginToDelete
	externalLoginStorageMock.LoginsByClientId[loginToDelete.ClientId] = loginToDelete

	err := externalLoginService.DeleteExternalLogin(loginToDelete.Id)
	assert.NoError(t, err)
	_, exists := externalLoginStorageMock.LoginsById[loginToDelete.Id]
	assert.False(t, exists)

	// Test deleting non-existent
	err = externalLoginService.DeleteExternalLogin(999)
	assert.Error(t, err)

	externalLoginStorageMock.DeleteErr = errors.New("forced Delete error")
	externalLoginStorageMock.LoginsById[loginToDelete.Id] = loginToDelete
	externalLoginStorageMock.LoginsByClientId[loginToDelete.ClientId] = loginToDelete
	err = externalLoginService.DeleteExternalLogin(loginToDelete.Id)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Delete error")
}
