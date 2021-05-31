package repository

import (
	"DBproject/internal/posts"
	"DBproject/models"
	"database/sql"
	"time"
)

type postsRepo struct {
	DB *sql.DB
}

func NewPostsRepo(db *sql.DB) posts.PostsRepo {
	return &postsRepo{
		DB: db,
	}
}

func (db *postsRepo) GetPost() (models.Post, models.Error) {
	return models.Post{}, models.Error{} // todo
}

func (db *postsRepo) UpdatePost(id int, message string) (models.Post, models.Error) {
	post := models.Post{
		ID: id,
	}
	err := db.DB.QueryRow("update posts set message = $1, isedited = true where id = $2 " +
		"returning parent, author, forum, thread, created", message, id).Scan(
			post.Parent, post.Author, post.Forum, post.Thread, post.Created)
	post.IsEdited = true
	post.Message = message

	if err == sql.ErrNoRows {
		return models.Post{}, models.Error{Code: 404}
	}

	return post, models.Error{}
}

func (db *postsRepo) CreatePosts(thread models.Thread, posts []models.Post) ([]models.Post, models.Error) {
	createdTime := time.Now()

	query := `
	insert into posts (parent, author, message, isedited, forum, thread, created) 
	values ($1,$2,$3,$4,$5,$6,$7) returning id, created
`

	for _, post := range posts {
		// todo isedited?
		err := db.DB.QueryRow(query, post.Parent, post.Author, post.Message, false, thread.Slug, thread.Forum, createdTime).
			Scan(&post.ID, &post.Created)
		if err != nil {
			return nil, models.Error{Code: 409}
		}
	}

	return posts, models.Error{}
}

func (db *postsRepo) GetThreadAndForumById(id int) (models.Thread, models.Error) {
	var thread models.Thread
	err := db.DB.QueryRow("select slug, forum from threads where id=$1", id).
		Scan(&thread.Slug, &thread.Forum)
	if err != nil {
		return models.Thread{}, models.Error{Code: 404}
	}

	return thread, models.Error{}
}

func (db *postsRepo) GetThreadAndForumBySlug(slug string) (models.Thread, models.Error) {
	var thread models.Thread
	err := db.DB.QueryRow("select id, forum from threads where slug=$1", slug).
		Scan(&thread.ID, &thread.Forum)
	if err != nil {
		return models.Thread{}, models.Error{Code: 404}
	}

	return thread, models.Error{}
}
