package repository

import (
	"DBproject/internal/posts"
	"DBproject/models"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"time"
)

type postsRepo struct {
	DB *pgx.ConnPool
	//	DbCreate *pgxpool.Pool
}

func NewPostsRepo(db *pgx.ConnPool) posts.PostsRepo {
	return &postsRepo{
		DB: db,
		//	DbCreate: DbCreate,
	}
}

func (db *postsRepo) GetPost(id int) (models.Post, models.Error) {
	post := models.Post{}
	var parent *int
	err := db.DB.QueryRow("select parent, author, message, isedited, forum, thread, created from posts where id=$1", id).
		Scan(&parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
	fmt.Println("get post error", err)
	if err == pgx.ErrNoRows {
		return models.Post{}, models.Error{Code: 404}
	}

	post.ID = id
	if parent != nil {
		post.Parent = *parent
	}

	return post, models.Error{Code: 200}
}

func (db *postsRepo) GetPostAuthor(pid int) (models.User, models.Error) {
	author := models.User{}
	err := db.DB.QueryRow(`
	select nickname, fullname, about, email from users
	inner join posts p on users.nickname = p.author
	where p.id = $1`, pid).Scan(&author.Nickname, &author.FullName, &author.About, &author.Email)
	fmt.Println("get post author error", err)
	if err != nil {
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
	fmt.Println("get post thread error", err)
	if err != nil {
		return thread, models.Error{Code: 500}
	}

	return thread, models.Error{Code: 200}
}

func (db *postsRepo) GetPostForum(pid int) (models.Forum, models.Error) {
	forum := models.Forum{}
	err := db.DB.QueryRow(`
	select f.title, f.user, f.slug, f.posts_count, f.threads_count from forums f
	join posts p on f.slug = p.forum
	where p.id=$1`, pid).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	fmt.Println("get post forum error", err)
	if err != nil {
		return forum, models.Error{Code: 500}
	}

	return forum, models.Error{Code: 200}
}

func (db *postsRepo) UpdatePost(id int, message string) (models.Post, models.Error) {
	post := models.Post{
		ID: id,
	}
	var parent *int
	post.IsEdited = true
	err := db.DB.QueryRow("update posts set message=coalesce(nullif($1, ''), message), "+
		"isedited = case when message=$1 or $1='' then isEdited else true end "+
		"where id = $2 "+
		"returning parent, author, forum, thread, created, message, isEdited", message, id).Scan(
		&parent, &post.Author, &post.Forum, &post.Thread, &post.Created, &post.Message, &post.IsEdited)
	fmt.Println("update post error", err)
	if parent != nil {
		post.Parent = *parent
	}

	if err != nil {
		return models.Post{}, models.Error{Code: 404}
	}
	if err == pgx.ErrNoRows {
		return models.Post{}, models.Error{Code: 404}
	}

	return post, models.Error{}
}

func (db *postsRepo) CreatePosts(thread models.Thread, posts []models.Post) ([]models.Post, models.Error) {
	if len(posts) != 0 && posts[0].Parent != 0 {
		var parentCheck int
		err := db.DB.QueryRow("select thread from posts where id = $1", posts[0].Parent).Scan(&parentCheck)
		if err != nil {
			return nil, models.Error{Code: 409, Message: "Parent post was created in another thread"}
		}
	}

	createdTime := time.Now()

	query := `insert into posts (parent, author, message, forum, thread, created) values `
	queryParams := make([]interface{}, 0)

	last := len(posts) - 1
	for i, post := range posts {
		fmt.Println("post", post)

		if !db.CheckValidParent(thread.ID, post.Parent) {
			return nil, models.Error{Code: 409}
		}

		var parentPost *int
		if post.Parent != 0 {
			parentPost = new(int)
			*parentPost = post.Parent
		}

		if i == last {
			query += fmt.Sprintf(`(nullif($%d,0),$%d,$%d,$%d,$%d,$%d) `, i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		} else {
			//query += fmt.Sprintf(`(%d,'%s','%s','%s', %d, $1), `, parentPost, post.Author, post.Message, thread.Forum, thread.ID)
			query += fmt.Sprintf(`(nullif($%d,0),$%d,$%d,$%d,$%d,$%d), `, i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		}
		queryParams = append(queryParams, parentPost, post.Author, post.Message, thread.Forum, thread.ID, createdTime)
	}

	query += " returning id, created"

	//transaction, err := db.DB.Begin()
	//rows, err := transaction.Query(query, queryParams...)
	//if err != nil {
	//	transaction.Rollback()
	//	_, ok := err.(pgx.PgError)
	//	if ok {
	//		return nil, models.Error{Code: 404, Message: fmt.Sprintf("%d", 1)}
	//	}
	//}

	//transaction, err := db.DbCreate.Begin(context.Background())
	//batch := new(pgx.Batch)
	//batch.Queue(query, queryParams...)
	rows, err := db.DB.Query(query, queryParams...)

	fmt.Println("ADDING POST ERROR ", err)
	if err != nil {
		return nil, models.Error{Code: 500}
	}

	i := 0
	for rows.Next() {
		err = rows.Scan(&posts[i].ID, &posts[i].Created)
		if err != nil {
			return nil, models.Error{Code: 500}
		}
		posts[i].Forum = thread.Forum
		posts[i].Thread = thread.ID
		fmt.Println(posts[i])
		i++
	}

	if dbErr, ok := rows.Err().(pgx.PgError); ok {
		fmt.Println("PGX ERROR")
		switch dbErr.Code {
		case pgerrcode.RaiseException:
			fmt.Println("40404")
			return nil, models.Error{Code: 404, Message: "Post parent not found"}
		case "23503":
			fmt.Println("23503")
			return nil, models.Error{Code: 404, Message: "User not found"}
		}
	}

	return posts, models.Error{}
}

func (db *postsRepo) GetThreadAndForumById(id int) (models.Thread, models.Error) {
	fmt.Println("      GetThreadAndForumById")

	var thread models.Thread
	err := db.DB.QueryRow("select slug, forum from threads where id=$1", id).
		Scan(&thread.Slug, &thread.Forum)
	thread.ID = id
	fmt.Println(thread)
	if err != nil || thread.Forum == "" {
		return models.Thread{}, models.Error{Code: 404, Message: "Can't find post thread by id:"}
	}

	return thread, models.Error{}
}

func (db *postsRepo) GetThreadAndForumBySlug(slug string) (models.Thread, models.Error) {
	fmt.Println("      GetThreadAndForumBySlug")

	var thread models.Thread
	err := db.DB.QueryRow("select id, forum from threads where slug=$1", slug).
		Scan(&thread.ID, &thread.Forum)
	thread.Slug = slug
	fmt.Println(thread)
	fmt.Println(thread.ID, thread.Forum)
	if err != nil || thread.ID == 0 {
		fmt.Println("ERRROROROR")
		return models.Thread{}, models.Error{Code: 404, Message: "Can't find post thread by id:"}
	}

	return thread, models.Error{}
}

func (db *postsRepo) CheckValidParent(thread, parent int) bool {
	if parent == 0 {
		return true
	}

	post, err := db.GetPost(parent)
	if err.Code != 200 {
		fmt.Println(11111111)
		return false
	}

	if post.ID == 0 {
		fmt.Println(22222)
		return false
	}

	if post.Thread != thread {
		fmt.Println(33333333)
		fmt.Println("post thread", post.Thread, "thread", thread)
		return false
	}

	return true
}
