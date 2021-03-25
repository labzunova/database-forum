package repository

import (
	"DBproject/internal/forum"
	"DBproject/internal/forum/delivery/http"
	"DBproject/models"
	"database/sql"
)

type forumRepo struct {
	DB *sql.DB
}

func NewForumRepo(db *sql.DB) forum.ForumRepo {
	return &forumRepo{
		DB: db,
	}
}

func (db *forumRepo) CreateNewForum(forum models.Forum) (models.Forum, models.Error) {
	return models.Forum{}, nil // todo
}

func (db *forumRepo) GetForum(id string) (models.Forum, models.Error) {
	return models.Forum{}, nil // todo
}

func (db *forumRepo) GetUsers(slug string, params http.UsersParse) ([]models.User, models.Error) {
	return []models.User{}, nil // todo
}

func (db *forumRepo) GetThreads(slug string, params http.UsersParse) ([]models.Thread, models.Error) {
	return []models.Thread{}, nil // todo
}
