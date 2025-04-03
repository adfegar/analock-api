package services

import (
	"github.com/adfer-dev/analock-api/models"
	"github.com/adfer-dev/analock-api/storage"
)

type UserBody struct {
	Email    string `json:"email" validate:"required,email"`
	UserName string `json:"username" validate:"required,alphanum"`
}

var userStorage *storage.UserStorage = &storage.UserStorage{}

func GetUserById(id uint) (*models.User, error) {
	user, err := userStorage.Get(id)

	if err != nil {
		return nil, err
	}

	return user.(*models.User), nil
}

func GetUserByEmail(email string) (*models.User, error) {
	user, err := userStorage.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	return user.(*models.User), nil
}

func SaveUser(userBody UserBody) (*models.User, error) {
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

func UpdateUser(userBody UserBody) (*models.User, error) {
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

func DeleteUser(id uint) error {
	return userStorage.Delete(id)
}
