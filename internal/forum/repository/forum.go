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
	err := f.DB.QueryRow("select title from forums where user=$1", forum.User).Scan()
	if err == sql.ErrNoRows {
		return models.Forum{}, models.Error{Code: 404}
	}

	err = f.DB.QueryRow("insert into forums (title, user, slug) values ($1, $2, $3)",
		forum.Title, forum.User, forum.Slug).Scan()
	if err != nil {
		DBerror, _ := err.(pgx.PgError) // TODO error handling
		switch DBerror.Code {
		case pgerrcode.UniqueViolation: // если такой форум уже еть
			// todo вернуть этот форум
			return models.Forum{}, models.Error{Code: 409}
		case pgerrcode.NotNullViolation: // если владелец не найден ???
			return models.Forum{}, models.Error{Code: 404}
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

	return *forum, models.Error{}
}

// TODO limit offset desc
func (f *forumRepo) GetUsers(slug string, params models.ParseParams) ([]models.User, models.Error) {
	var forumID int
	err := f.DB.QueryRow("select id from forums where slug = $1", slug).Scan(forumID)
	if err != nil {

	}

	forumUsers, err := f.DB.Query("select f.nickname, f.fullname, f.email, f.about from" +
		"forum_users uf " +
		"join forums f where f.id = uf.forumID " +
		"join users u where uf.userID = u.id " +
		"where uf.forumID = $1 " +
		"order by uf.nickname", forumID)

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
			return  []models.User{}, models.Error{}
		}

		users = append(users, *user)
	}

	return users, models.Error{}
}

func (f *forumRepo) GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error) {
	var forumTitle int
	err := f.DB.QueryRow("select title from forums where slug = $1", slug).Scan(forumTitle)
	if err != nil {

	}

	forumUsers, err := f.DB.Query("select t.id, t.title, t.author, t.forum, t.message, t.votes, t.slug, t.created from " +
		"threads t where t.forum = $1" +
		"order by t.forum", forumTitle)

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
			return  []models.Thread{}, models.Error{}
		}

		threads = append(threads, *thread)
	}

	return threads, models.Error{}}
