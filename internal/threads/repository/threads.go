package repository

import (
	"DBproject/internal/threads"
	"DBproject/models"
	"database/sql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
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
	var forumID int
	err := db.DB.QueryRow("select id from forums where slug = $1", slug).Scan(forumID)
	if err != nil {

	}

	err = db.DB.QueryRow("insert into threads (title, author, message, created, forum, slug) values ($1,$2,$3,$4,$5) returning id",
		thread.Title, thread.Author, thread.Message, thread.Created, forumID, slug).Scan(thread.ID)
	thread.Votes = 0
	DBerror, _ := err.(pgx.PgError) // TODO error handling
	switch DBerror.Code {
	case pgerrcode.UniqueViolation: // если такой тред уже еть
		// todo вернуть этот тред
		return models.Thread{}, models.Error{Code: 409}
	case pgerrcode.NotNullViolation: // если владелец не найден ???
		return models.Thread{}, models.Error{Code: 404}
	}

	return models.Thread{}, models.Error{} // todo
}

func (db *threadsRepo) GetThread(slug string) (models.Thread, models.Error) {
	return models.Thread{}, models.Error{} // todo
}

func (db *threadsRepo) UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	return models.Thread{}, models.Error{} // todo
}

func (db *threadsRepo) GetThreadPosts(slug string, params models.ParseParamsThread) ([]models.Post, models.Error) {
	return []models.Post{}, models.Error{} // todo
}

func (db *threadsRepo) VoteThread(slug string, vote models.Vote) (models.Thread, models.Error) {
	return models.Thread{}, models.Error{} // todo
}
