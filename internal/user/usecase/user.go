package usecase

import (
	"DBproject/internal/user"
	"DBproject/models"
)

type userUsecase struct {
	userRepository user.UserRepo
}

func NewUserUsecase(repo user.UserRepo) user.UserUsecase {
	return &userUsecase{
		userRepository: repo,
	}
}

func (u *userUsecase) Create(user models.User) (models.User, models.Error) {
	return u.userRepository.CreateUser(user)
}

func (u *userUsecase) GetByNickname(nickname string) (models.User, models.Error) {
	return u.userRepository.GetUser(nickname)
}

func (u *userUsecase) Update(profle models.User) (models.User, models.Error) {
	return u.userRepository.UpdateUser(profle)
}
