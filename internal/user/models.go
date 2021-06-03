package user

import (
	"DBproject/models"
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	UserCreate(c echo.Context) error
	UserGetOne(c echo.Context) error
	UserUpdate(c echo.Context) error
}

type UserUsecase interface {
	Create(user models.User) models.Error
	GetByNickname(nickname string) (models.User, models.Error)
	Update(profile models.User) (models.User, models.Error)
	GetExistingUsers(nickname, email string) ([]models.User, models.Error)
}

type UserRepo interface {
	CreateUser(profile models.User) models.Error
	GetUser(nickname string) (models.User, models.Error)
	UpdateUser(profile models.User) (models.User, models.Error)
	GetExistingUsers(nickname, email string) ([]models.User, models.Error)
}