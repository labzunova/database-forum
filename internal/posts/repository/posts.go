package repository

import (
	"DBproject/internal/posts"
	"DBproject/models"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"time"
)

type postsRepo struct {
	DB *pgx.ConnPool
}

func NewPostsRepo(db *pgx.ConnPool) posts.PostsRepo {
	return &postsRepo{
		DB: db,
	}
}

func (db *postsRepo) GetPost(id int) (models.Post, models.Error) {
	post := models.Post{}
	err := db.DB.QueryRow("select parent, author, message, isedited, forum, thread, created from posts where id=$1", id).
		Scan(&post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	if err == sql.ErrNoRows {
		return models.Post{}, models.Error{Code: 404}
	}

	return post, models.Error{Code: 200}
}

func (db *postsRepo) GetPostAuthor(pid int) (models.User, models.Error) {
	author := models.User{}
	err := db.DB.QueryRow(`
	select nickname, fullname, about, email from users
	inner join posts p on users.nickname = p.author
	where p.id = $1`, pid).Scan(&author.Nickname, &author.FullName, &author.About, &author.Email)
	if err!= nil {
		return author, models.Error{Code: 500}
	}

	return author, models.Error{Code: 200}
}

func (db *postsRepo) GetPostThread(pid int) (models.Thread, models.Error) {
	thread := models.Thread{}
	err := db.DB.QueryRow(`
	select t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from threads t
	join posts p on t.id = p.thread
	where p.id = $1`, pid).Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		return thread, models.Error{Code: 500}
	}

	return thread, models.Error{Code: 200}
}

func (db *postsRepo) GetPostForum(pid int) (models.Forum, models.Error) {
	forum := models.Forum{}
	err := db.DB.QueryRow(`
	select f.title, f.user, f.slug, f.posts, f.threads from forums f
	join posts p on f.slug = p.forum
	where p.id=$1`, pid).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return forum, models.Error{Code: 500}
	}

	return forum, models.Error{Code: 200}
}

func (db *postsRepo) UpdatePost(id int, message string) (models.Post, models.Error) {
	post := models.Post{
		ID: id,
	}
	err := db.DB.QueryRow("update posts set message = $1, isedited = true where id = $2 " +
		"returning parent, author, forum, thread, created", message, id).Scan(
			&post.Parent, &post.Author, &post.Forum, &post.Thread, &post.Created)
	post.IsEdited = true
	post.Message = message

	if err == sql.ErrNoRows {
		return models.Post{}, models.Error{Code: 404}
	}

	return post, models.Error{}
}

func (db *postsRepo) CreatePosts(thread models.Thread, posts []models.Post) ([]models.Post, models.Error) {
	createdTime := time.Now()

	query := `insert into posts (parent, author, message, forum, thread, created) values `
	queryParams := make([]interface{}, 0)

	last := len(posts) - 1
	for i, post := range posts {
		fmt.Println("post", post)
		//parentPost := new(int)
		//if post.Parent != 0 {
		//	*parentPost = post.Parent
		//}
		//var parentPost *int
		//if post.Parent != 0 {
		//	parentPost = new(int)
		//	*parentPost = post.Parent
		//}

	//	fmt.Println(parentPost)
		if i == last {
			query += fmt.Sprintf(`(nullif($%d,0),$%d,$%d,$%d,$%d,$%d) `, i*6+1,i*6+2,i*6+3,i*6+4,i*6+5,i*6+6)
		} else {
			//query += fmt.Sprintf(`(%d,'%s','%s','%s', %d, $1), `, parentPost, post.Author, post.Message, thread.Forum, thread.ID)
			query += fmt.Sprintf(`(nullif($%d,0),$%d,$%d,$%d,$%d,$%d), `, i*6+1,i*6+2,i*6+3,i*6+4,i*6+5,i*6+6)
		}
		queryParams = append(queryParams, post.Parent, post.Author, post.Message, thread.Forum, thread.ID, createdTime)
	}

	query += " returning id, created"

	fmt.Println("\n", query, "\n")
	rows, err := db.DB.Query(query, queryParams...)
	fmt.Println("ADDING POST ERROR ", err)
	if err != nil {
		return nil, models.Error{Code: 500}
	}

	i := 0
	for rows.Next() {
		fmt.Println("scan iteration")
		err = rows.Scan(&posts[i].ID, &posts[i].Created)
		if err != nil {
			return nil, models.Error{Code: 500}
		}
		posts[i].Forum = thread.Forum
		posts[i].Thread = thread.ID
		fmt.Println(posts[i])
		i++
	}
fmt.Println(posts)
	return posts, models.Error{}
}

func (db *postsRepo) GetThreadAndForumById(id int) (models.Thread, models.Error) {
	fmt.Println("      GetThreadAndForumById")

	var thread models.Thread
	err := db.DB.QueryRow("select slug, forum from threads where id=$1", id).
		Scan(&thread.Slug, &thread.Forum)
	thread.ID = id
	fmt.Println(thread)
	if err != nil {
		return models.Thread{}, models.Error{Code: 404}
	}

	return thread, models.Error{}
}

func (db *postsRepo) GetThreadAndForumBySlug(slug string) (models.Thread, models.Error) {
	fmt.Println("      GetThreadAndForumBySlug")

	var thread models.Thread
	err := db.DB.QueryRow("select id, forum from threads where slug=$1", slug).
		Scan(&thread.ID, &thread.Forum)
	thread.Slug = slug
	if err != nil {
		return models.Thread{}, models.Error{Code: 404}
	}

	return thread, models.Error{}
}
