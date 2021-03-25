package repository

import (
	"DBproject/internal/threads"
	"DBproject/internal/threads/delivery/http"
	"DBproject/models"
	"database/sql"
)

type threadsRepo struct {
	DB *sql.DB
}

func NewThreadsRepo(db *sql.DB) threads.ThreadsRepo {
	return &threadsRepo{
		DB: db,
	}
}

func (db *threadsRepo) CreateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	return models.Thread{}, nil // todo

}

func (db *threadsRepo) GetThread(slug string) (models.Thread, models.Error) {
	return models.Thread{}, nil // todo
}

func (db *threadsRepo) UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	return models.Thread{}, nil // todo
}

func (db *threadsRepo) GetThreadPosts(slug string, params http.ThreadsParse) ([]models.Post, models.Error) {
	return []models.Post{}, nil // todo
}

func (db *threadsRepo) VoteThread(slug string, vote models.Vote) (models.Thread, models.Error) {
	return models.Thread{}, nil // todo
}
