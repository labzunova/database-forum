package repository

import (
	"DBproject/internal/forum"
	"DBproject/models"
	"database/sql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type forumRepo struct {
	DB *sql.DB
}

func NewForumRepo(db *sql.DB) forum.ForumRepo {
	return &forumRepo{
		DB: db,
	}
}

func (f *forumRepo) CreateNewForum(forum models.Forum) (models.Forum, models.Error) {
	// checkUserExists
	err := f.DB.QueryRow("insert into forums (title, user, slug) values ($1, $2, $3)",
		forum.Title, forum.User, forum.Slug).Scan()
	if err != nil {
		DBerror, _ := err.(pgx.PgError) // TODO error handling
		switch DBerror.Code {
		case pgerrcode.UniqueViolation: // если такой форум уже еть
			// todo вернуть этот форум
			return models.Forum{}, models.Error{Code: 409}
		case pgerrcode.NotNullViolation: // если владелец не найден ???
			return models.Forum{}, models.Error{Code: 404}
		default:
			return models.Forum{}, models.Error{Code: 500}
		}
	}

	return forum, models.Error{}
}

func (f *forumRepo) GetForum(slug string) (models.Forum, models.Error) {
	forum := new(models.Forum)
	err := f.DB.QueryRow("select title, user, slug, posts, threads from forums where slug = $1",
		slug).Scan(forum.Title, forum.User, forum.Slug, forum.Posts, forum.Threads)
	if err != nil {
		return models.Forum{}, models.Error{Code: 404}
	}

	return *forum, models.Error{Code: 200}
}

func (f *forumRepo) CreateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	err := f.DB.QueryRow(
		`
	insert into threads 
    (title, author, message, forum, slug) 
	values ($2,$3,$4,$5) 
	returning id, u.nickname, created`,
		thread.Title, thread.Author, thread.Message, thread.Forum, slug).Scan(thread.ID, thread.Author)

	DBerror, _ := err.(pgx.PgError) // TODO error handling
	switch DBerror.Code {
	case pgerrcode.UniqueViolation: // если такой тред уже еть
		return thread, models.Error{Code: 409}
	case pgerrcode.NotNullViolation: // если владелец не найден ???
		return models.Thread{}, models.Error{Code: 404}
	}

	return thread, models.Error{Code: 200}// todo
}

func (f *forumRepo) GetUsers(slug string, params models.ParseParams) ([]models.User, models.Error) {
	var queryParametres []interface{}
	query := `
		select u.nickname, u.fullname, u.email, u.about from
		forum_users uf 
		join users u where uf.userID = u.id 
		where uf.forumSlug = $1
	`
	queryParametres = append(queryParametres, slug)

	if params.Since != "" {
		query += " and nickname > $2 "
		queryParametres = append(queryParametres, params.Since)
	}

	if !params.Desc {
		query += " order by uf.nickname"
	} else {
		query += " order by uf.nickname desc"
	}

	if params.Limit != 0 {
		if params.Since == "" {
			query += "  LIMIT $2"
		} else {
			query += "  LIMIT $3"
		}
		queryParametres = append(queryParametres, params.Limit)
	}

	forumUsers, err := f.DB.Query(query, queryParametres)
	if err != nil {
		return []models.User{}, models.Error{Code: 404}
	}

	users := make([]models.User, 0)
	for forumUsers.Next() {
		user := new(models.User)
		err = forumUsers.Scan(
			user.Nickname,
			user.FullName,
			user.Email,
			user.About,
		)
		if err != nil {
			return []models.User{}, models.Error{Code: 500}
		}

		users = append(users, *user)
	}

	return users, models.Error{}
}

func (f *forumRepo) GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error) {
	var queryParams []interface{}
	query := `
		select t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from
		threads t where t.slug = $1 
	`
	queryParams = append(queryParams, slug)

	if params.Since != "" {
		query += " and created > $2 "
		queryParams = append(queryParams, params.Since)
	}

	if !params.Desc {
		query += " order by t.forum"
	} else {
		query += " order by t.forum desc"
	}

	if params.Limit != 0 {
		if params.Since == "" {
			query += "  LIMIT $2"
		} else {
			query += "  LIMIT $3"
		}
		queryParams = append(queryParams, params.Limit)
	}

	forumUsers, err := f.DB.Query(query, queryParams)
	if err != nil {
		return []models.Thread{}, models.Error{Code: 404}
	}

	threads := make([]models.Thread, 0)
	for forumUsers.Next() {
		thread := new(models.Thread)
		err = forumUsers.Scan(
			thread.ID,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Votes,
			thread.Slug,
			thread.Created,
		)
		if err != nil {
			return  []models.Thread{}, models.Error{Code: 500}
		}

		threads = append(threads, *thread)
	}

	return threads, models.Error{Code: 200}
}
