package repository

import (
	"DBproject/internal/user"
	"DBproject/models"
	"database/sql"
	"fmt"
)

type usersRepo struct {
	DB *sql.DB
}

func NewUsersRepo(db *sql.DB) user.UserRepo {
	return &usersRepo{
		DB: db,
	}
}

func (db *usersRepo) CreateUser(profile models.User) models.Error {
	fmt.Println("create user", profile)

	_, err := db.DB.Exec(`insert into users (nickname, fullname, about, email) values ($1,$2,$3,$4)`,
		profile.Nickname, profile.FullName, profile.About, profile.Email)
	//dbError, ok := err.(pgx.PgError)
	//if ok && dbError.Code == pgerrcode.UniqueViolation {
	//	return models.Error{Code: 409}
	//}
	if err != nil {
		return models.Error{Code: 409}
	}

	return models.Error{Code: 201}
}

func (db *usersRepo) GetUser(nickname string) (models.User, models.Error) {
	fmt.Println("get user ", nickname)

	user := models.User{}
	err := db.DB.QueryRow("select nickname, fullname, about, email from users where nickname=$1", nickname).
		Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)
	if err == sql.ErrNoRows {
		return models.User{}, models.Error{Code: 404}
	}
	if err != nil {
		return models.User{}, models.Error{Code: 500}
	}

	return user, models.Error{Code: 200}
}

func (db *usersRepo) UpdateUser(profile models.User) (models.User, models.Error) {
	fmt.Println("update user", profile)

	err := db.DB.QueryRow(`
		update users set 
		fullname=coalesce(nullif($1, ''), fullname),
		about=coalesce(nullif($2, ''), about),
		email=coalesce(nullif($3, ''), email)
		where nickname=$4
		returning fullname, about, email`, profile.FullName, profile.About, profile.Email, profile.Nickname).
		Scan(&profile.FullName, &profile.About, &profile.Email)
	//dbError, ok := err.(pgx.PgError)
	if err == sql.ErrNoRows {
		return models.User{}, models.Error{Code: 404}
	}
	//if ok && dbError.Code == pgerrcode.UniqueViolation {
	//	return models.User{}, models.Error{Code: 409}
	//}
	if err != nil {
		return models.User{}, models.Error{Code: 409}
	}

	return profile, models.Error{Code: 200}
}

// extra

func (db *usersRepo) GetExistingUsers(nickname, email string) ([]models.User, models.Error) {
	users := make([]models.User, 0)

	rows, err := db.DB.Query(`select nickname, fullname, about, email from users where nickname=$1 or email=$2`,
		nickname, email)
	if err != nil {
		return []models.User{}, models.Error{Code: 500}
	}

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
			&user.Nickname,
			&user.FullName,
			&user.About,
			&user.Email)
		if err != nil {
			return []models.User{}, models.Error{Code: 500}
		}

		users = append(users, user)
	}

	return users, models.Error{Code: 200}
}