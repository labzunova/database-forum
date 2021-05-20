package repository

import (
	"DBproject/internal/posts"
	"DBproject/models"
	"database/sql"
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

func (db *postsRepo) UpdatePost(id int64, message string) (models.Post, models.Error) {
	post := new(models.Post)
	err := db.DB.QueryRow("update posts set message = $1, isedited = true where id = $2 " +
		"returning parent, author, forum, thread, created", message, id).Scan(
			post.Parent, post.Author, post.Forum, post.Thread, post.Created)
	post.IsEdited = true
	post.Message = message

	if err == sql.ErrNoRows {
		return models.Post{}, models.Error{Code: 404}
	}

	return *post, models.Error{}
}

func (db *postsRepo) CreatePosts(slug string, posts []models.Post) ([]models.Post, models.Error) {
	newPosts := make([]models.Post, 0)
	var forumName string
	err := db.DB.QueryRow("select title from forums where slug = $1", slug).Scan(forumName)
	if err != nil {
		return newPosts, models.Error{Code: 404}
	}

	for _, post := range posts {
		postsDB, err := db.DB.QueryRow("insert into posts (parent, author, message, isedited, forum, thread, created) values ($1,$2,$3,$4,$5,$6,$7)")
		// todo check
	}


	return []models.Post{}, models.Error{} // todo
}
