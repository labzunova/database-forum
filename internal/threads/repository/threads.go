package repository

import (
	"DBproject/internal/threads"
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

func (db *threadsRepo) GetThread(slug string, id int) (models.Thread, models.Error) {
	thread := models.Thread{
		Slug: slug,
		ID: id,
	}
	if id == 0 {
		err := db.DB.QueryRow("select title, author, forum, message, votes, slug, created from threads where id = $1", id).
			Scan(&thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return thread, models.Error{Code: 404}
		}
	} else {
		err := db.DB.QueryRow("select id, title, author, forum, message, votes, created from threads where slug = $1", slug).
			Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created)
		if err != nil {
			return thread, models.Error{Code: 404}
		}
	}

	return thread, models.Error{}
}

func (db *threadsRepo) UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	err := db.DB.QueryRow(`
		update thread set title=$1, message=$2
		where slug = $3 
		returning id, author, forum, votes, created `,
		thread.Title, thread.Message, slug).
		Scan()
	if err != nil {

	}

}

func (db *threadsRepo) GetThreadPosts(slug string, params models.ParseParamsThread) ([]models.Post, models.Error) {
	return []models.Post{}, models.Error{} // todo
}

func (db *threadsRepo) VoteThread(slug string, vote models.Vote) (models.Thread, models.Error) {
	return models.Thread{}, models.Error{} // todo
}
