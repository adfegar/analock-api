package services

import (
	"errors"
	"fmt"
	"testing"

	"github.com/adfer-dev/analock-api/models"
	"github.com/stretchr/testify/assert"
)

// userStorageMockUserStorage implements UserStorageInterface
type userStorageMockUserStorage struct {
	UsersByEmail  map[string]*models.User
	UsersById     map[uint]*models.User
	GetErr        error
	GetByEmailErr error
	CreateErr     error
	UpdateErr     error
	DeleteErr     error
	nextId        uint
}

func newuserStorageMockUserStorage() *userStorageMockUserStorage {
	return &userStorageMockUserStorage{
		UsersByEmail: make(map[string]*models.User),
		UsersById:    make(map[uint]*models.User),
		nextId:       1,
	}
}

func (m *userStorageMockUserStorage) Get(id uint) (interface{}, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	user, ok := m.UsersById[id]
	if !ok {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

func (m *userStorageMockUserStorage) GetByEmail(email string) (interface{}, error) {
	if m.GetByEmailErr != nil {
		return nil, m.GetByEmailErr
	}
	user, ok := m.UsersByEmail[email]
	if !ok {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	return user, nil
}

func (m *userStorageMockUserStorage) Create(data interface{}) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	user, ok := data.(*models.User)
	if !ok {
		return errors.New("create: invalid type for User")
	}
	if user.Id == 0 {
		user.Id = m.nextId
		m.nextId++
	}
	m.UsersById[user.Id] = user
	m.UsersByEmail[user.Email] = user
	return nil
}

func (m *userStorageMockUserStorage) Update(data interface{}) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	user, ok := data.(*models.User)
	if !ok {
		return errors.New("update: invalid type for User")
	}

	// Try finding by ID first if it's non-zero
	if user.Id != 0 {
		if _, ok := m.UsersById[user.Id]; ok {
			m.UsersById[user.Id] = user
			m.UsersByEmail[user.Email] = user
			return nil
		}
	}

	// If not found by ID, or ID is zero, try to find by email
	var existingUserById *models.User
	for _, u := range m.UsersByEmail {
		if u.Email == user.Email {
			existingUserById = u
			break
		}
	}

	if existingUserById != nil {
		// Update the fields
		existingUserById.UserName = user.UserName
		m.UsersByEmail[user.Email] = existingUserById
		m.UsersById[existingUserById.Id] = existingUserById
		return nil
	}

	return fmt.Errorf("update: user with email %s not found to update", user.Email)
}

func (m *userStorageMockUserStorage) Delete(id uint) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	user, ok := m.UsersById[id]
	if !ok {
		return fmt.Errorf("delete: user with id %d not found", id)
	}
	delete(m.UsersByEmail, user.Email)
	delete(m.UsersById, id)
	return nil
}

var userService UserService = &UserServiceImpl{}

func TestGetUserById(t *testing.T) {
	originalStorage := userStorage
	userStorageMock := newuserStorageMockUserStorage()
	userStorage = userStorageMock
	defer func() { userStorage = originalStorage }()

	testUser := &models.User{Id: 1, Email: "test@example.com", UserName: "testuser"}
	userStorageMock.UsersById[testUser.Id] = testUser
	userStorageMock.UsersByEmail[testUser.Email] = testUser

	user, err := userService.GetUserById(1)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user)

	_, err = userService.GetUserById(2)
	assert.Error(t, err)

	userStorageMock.GetErr = errors.New("forced Get error")
	_, err = userService.GetUserById(1)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Get error")
}

func TestGetUserByEmail(t *testing.T) {
	originalStorage := userStorage
	userStorageMock := newuserStorageMockUserStorage()
	userStorage = userStorageMock
	defer func() { userStorage = originalStorage }()

	testUser := &models.User{Id: 1, Email: "test@example.com", UserName: "testuser"}
	userStorageMock.UsersByEmail[testUser.Email] = testUser
	userStorageMock.UsersById[testUser.Id] = testUser

	user, err := userService.GetUserByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, testUser, user)

	_, err = userService.GetUserByEmail("nonexistent@example.com")
	assert.Error(t, err)

	userStorageMock.GetByEmailErr = errors.New("forced GetByEmail error")
	_, err = userService.GetUserByEmail("test@example.com")
	assert.Error(t, err)
	assert.EqualError(t, err, "forced GetByEmail error")
}

func TestSaveUser(t *testing.T) {
	originalStorage := userStorage
	userStorageMock := newuserStorageMockUserStorage()
	userStorage = userStorageMock
	defer func() { userStorage = originalStorage }()

	userBody := UserBody{Email: "new@example.com", UserName: "newuser"}

	// Test successful save
	savedUser, err := userService.SaveUser(userBody)
	assert.NoError(t, err)
	assert.NotNil(t, savedUser)
	assert.Equal(t, userBody.Email, savedUser.Email)
	assert.Equal(t, userBody.UserName, savedUser.UserName)
	assert.Equal(t, models.Standard, savedUser.Role)
	assert.True(t, savedUser.Id > 0) // check that storage assigned an ID

	// Verify it's in the userStorageMock storage
	_, okId := userStorageMock.UsersById[savedUser.Id]
	_, okEmail := userStorageMock.UsersByEmail[savedUser.Email]
	assert.True(t, okId)
	assert.True(t, okEmail)

	// Test error from storage.Create
	userStorageMock.CreateErr = errors.New("forced Create error")
	_, err = userService.SaveUser(UserBody{Email: "error@example.com", UserName: "erroruser"})
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Create error")
}

func TestUpdateUser(t *testing.T) {
	originalStorage := userStorage
	userStorageMock := newuserStorageMockUserStorage()
	userStorage = userStorageMock
	defer func() { userStorage = originalStorage }()

	// Pre-populate a user
	initialEmail := "update@example.com"
	initialUser := &models.User{Id: 5, Email: initialEmail, UserName: "initialuser", Role: models.Standard}
	userStorageMock.UsersById[initialUser.Id] = initialUser
	userStorageMock.UsersByEmail[initialUser.Email] = initialUser

	updateBody := UserBody{Email: initialEmail, UserName: "updateduser"}

	// Test successful update
	updatedUser, err := userService.UpdateUser(updateBody)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, updateBody.UserName, updatedUser.UserName)
	assert.Equal(t, updateBody.Email, updatedUser.Email)
	assert.Equal(t, models.Standard, updatedUser.Role)

	userInuserStorageMock := userStorageMock.UsersByEmail[initialEmail]
	assert.NotNil(t, userInuserStorageMock)
	assert.Equal(t, "updateduser", userInuserStorageMock.UserName)

	// Test error from storage.Update if user not found by email
	userStorageMock.UpdateErr = nil // reset
	_, err = userService.UpdateUser(UserBody{Email: "nonexistent@example.com", UserName: "ghost"})
	assert.Error(t, err)

	// Test forced error from storage.Update
	userStorageMock.UsersByEmail["forceerror@example.com"] = &models.User{Id: 6, Email: "forceerror@example.com", UserName: "pre"}
	userStorageMock.UsersById[6] = userStorageMock.UsersByEmail["forceerror@example.com"]
	userStorageMock.UpdateErr = errors.New("forced Update error")
	_, err = userService.UpdateUser(UserBody{Email: "forceerror@example.com", UserName: "forcingerror"})
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Update error")
}

func TestDeleteUser(t *testing.T) {
	originalStorage := userStorage
	userStorageMock := newuserStorageMockUserStorage()
	userStorage = userStorageMock
	defer func() { userStorage = originalStorage }()

	userToDelete := &models.User{Id: 10, Email: "delete@example.com", UserName: "deleteuser"}
	userStorageMock.UsersById[userToDelete.Id] = userToDelete
	userStorageMock.UsersByEmail[userToDelete.Email] = userToDelete

	// Test successful delete
	err := userService.DeleteUser(userToDelete.Id)
	assert.NoError(t, err)
	_, okId := userStorageMock.UsersById[userToDelete.Id]
	_, okEmail := userStorageMock.UsersByEmail[userToDelete.Email]
	assert.False(t, okId)
	assert.False(t, okEmail)

	// Test delete non-existent
	err = userService.DeleteUser(999) // ID that doesn't exist
	assert.Error(t, err)              // userStorageMock returns error for not found

	// Test forced error from storage.Delete
	userStorageMock.UsersById[userToDelete.Id] = userToDelete
	userStorageMock.UsersByEmail[userToDelete.Email] = userToDelete
	userStorageMock.DeleteErr = errors.New("forced Delete error")
	err = userService.DeleteUser(userToDelete.Id)
	assert.Error(t, err)
	assert.EqualError(t, err, "forced Delete error")
}
