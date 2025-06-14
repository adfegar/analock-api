package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

type UserBody struct {
	Email    string `json:"email" validate:"required,email"`
	UserName string `json:"username" validate:"required,alphanum"`
}

var userStorage storage.UserStorageInterface = &storage.UserStorage{}

// UserService defines all operations for the user service.
type UserService interface {
	GetUserById(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	SaveUser(userBody UserBody) (*models.User, error)
	UpdateUser(userBody UserBody) (*models.User, error)
	DeleteUser(id uint) error
}

// UserServiceImpl is the concrete implementation of UserService.
type UserServiceImpl struct{}

func (userService *UserServiceImpl) GetUserById(id uint) (*models.User, error) {
	user, err := userStorage.Get(id)
	if err != nil {
		return nil, err
	}
	return user.(*models.User), nil
}

func (userService *UserServiceImpl) GetUserByEmail(email string) (*models.User, error) {
	user, err := userStorage.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user.(*models.User), nil
}

func (userService *UserServiceImpl) SaveUser(userBody UserBody) (*models.User, error) {
	savedUser := &models.User{
		Email:    userBody.Email,
		UserName: userBody.UserName,
		Role:     models.Standard,
	}
	err := userStorage.Create(savedUser)
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}

func (userService *UserServiceImpl) UpdateUser(userBody UserBody) (*models.User, error) {
	updatedUser := &models.User{}
	updatedUser.UserName = userBody.UserName
	updatedUser.Email = userBody.Email
	updatedUser.Role = models.Standard
	err := userStorage.Update(updatedUser)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (userService *UserServiceImpl) DeleteUser(id uint) error {
	return userStorage.Delete(id)
}
