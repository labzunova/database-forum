package user

import (
	"DBproject/models"
	"net/http"
)

type UserHandler interface {
	UserCreate(w http.ResponseWriter, r *http.Request)
	UserGetOne(w http.ResponseWriter, r *http.Request)
	UserUpdate(w http.ResponseWriter, r *http.Request)
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
