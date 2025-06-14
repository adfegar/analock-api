package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

// UserStorageInterface defines storage operations for users.
type UserStorageInterface interface {
	Get(id uint) (interface{}, error)
	GetByEmail(email string) (interface{}, error)
	Create(data interface{}) error
	Update(data interface{}) error
	Delete(id uint) error
}

type UserBody struct {
	Email    string `json:"email" validate:"required,email"`
	UserName string `json:"username" validate:"required,alphanum"`
}

var userStorage UserStorageInterface = &storage.UserStorage{}

// UserService defines all operations for the user service.
type UserService interface {
	GetUserById(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	SaveUser(userBody UserBody) (*models.User, error)
	UpdateUser(userBody UserBody) (*models.User, error)
	DeleteUser(id uint) error
}

// DefaultUserService is the concrete implementation of UserService.
type DefaultUserService struct{}

func (s *DefaultUserService) GetUserById(id uint) (*models.User, error) {
	user, err := userStorage.Get(id)
	if err != nil {
		return nil, err
	}
	return user.(*models.User), nil
}

func (s *DefaultUserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := userStorage.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return user.(*models.User), nil
}

func (s *DefaultUserService) SaveUser(userBody UserBody) (*models.User, error) {
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

func (s *DefaultUserService) UpdateUser(userBody UserBody) (*models.User, error) {
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

func (s *DefaultUserService) DeleteUser(id uint) error {
	return userStorage.Delete(id)
}
