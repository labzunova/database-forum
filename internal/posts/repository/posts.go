package repository

import (
	"DBproject/models"
	"database/sql"
)

type PostsRepository interface {
	GetPost() ([]models.Post, models.Error)
	UpdatePost(id int64, post models.Post) (models.Post, models.Error)
	CreatePosts(slug string, posts []models.Post) ([]models.Post, models.Error)
}

type postsRepo struct {
	DB *sql.DB
}

func NewPostsRepo(db *sql.DB) PostsRepository {
	return &postsRepo{
		DB: db,
	}
}

func (db *postsRepo) GetPost() ([]models.Post, models.Error) {
	return []models.Post{}, nil // todo
}

func (db *postsRepo) UpdatePost(id int64, post models.Post) (models.Post, models.Error) {
	return models.Post{}, nil // todo
}

func (db *postsRepo) CreatePosts(slug string, posts []models.Post) ([]models.Post, models.Error) {
	return []models.Post{}, nil // todo
}
