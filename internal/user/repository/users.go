package repository

import (
	"DBproject/internal/user"
	"DBproject/models"
	"database/sql"
)

type usersRepo struct {
	DB *sql.DB
}

func NewUsersRepo(db *sql.DB) user.UserRepo {
	return &usersRepo{
		DB: db,
	}
}

func (db *usersRepo) CreateUser(profile models.User) (models.User, models.Error) {
	// TODO
}

func (db *usersRepo) GetUser(nickname string) (models.User, models.Error) {
	// TODO
}

func (db *usersRepo) UpdateUser(profile models.User) (models.User, models.Error) {

}
