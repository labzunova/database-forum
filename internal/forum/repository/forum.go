package repository

import (
	"DBproject/internal/forum"
	"DBproject/models"
	"database/sql"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type forumRepo struct {
	DB *pgx.ConnPool
}

func NewForumRepo(db *pgx.ConnPool) forum.ForumRepo {
	return &forumRepo{
		DB: db,
	}
}

func (f *forumRepo) CreateNewForum(forum models.Forum) (models.Forum, models.Error) {
	err := f.DB.QueryRow("select nickname from users where nickname=$1", forum.User).Scan(&forum.User)
	if err != nil {
		return models.Forum{}, models.Error{Code: 404}
	}

	// checkUserExists
	_, err = f.DB.Exec(`insert into forums 
    	(slug, title, "user") 
    	values ($1, $2, $3)`,
		forum.Slug, forum.Title, forum.User)
	fmt.Println("create err ", err)
	//dbErr, ok := err.(pgx.PgError)
	//fmt.Println("create err ", dbErr, ok, dbErr.Code, dbErr.Message)
	//
	//if ok {
	//		if dbErr.Code == pgerrcode.NotNullViolation || dbErr.Code == pgerrcode.ForeignKeyViolation { // если владелец не найден
	//		return models.Forum{}, models.Error{Code: 404}
	//	}
	//	if dbErr.Code == pgerrcode.UniqueViolation { // если такой форум уже еть
	//
	//}

	if err != nil && err != sql.ErrNoRows { // если такой форум уже еть
		fmt.Println("409")
		return models.Forum{}, models.Error{Code: 409}
	}
	//if err == sql.ErrNoRows { // если владелец не найден
	//	fmt.Println("404")
	//	return models.Forum{}, models.Error{Code: 404}
	//}

	return forum, models.Error{Code: 200}
}

func (f *forumRepo) GetForum(slug string) (models.Forum, models.Error) {
	forumm := new(models.Forum)
	err := f.DB.QueryRow(`select slug, title, threads_count, posts_count, "user" from forums where slug = $1`, // TODO limit 1?
		slug).Scan(&forumm.Slug, &forumm.Title, &forumm.Threads, &forumm.Posts, &forumm.User)
	fmt.Println("err", err)
	if err != nil {
		return models.Forum{}, models.Error{Code: 404}
	}

	return *forumm, models.Error{Code: 200}
}

func (f *forumRepo) CreateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	fmt.Println("create thread: ", slug, thread)

	var threadSlug *string
	if thread.Slug != "" {
		threadSlug = &thread.Slug
	}

	err := f.DB.QueryRow("select nickname from users where nickname=$1 limit 1", thread.Author).Scan(&thread.Author)
	if err == pgx.ErrNoRows {
		fmt.Println("404")
		return models.Thread{}, models.Error{Code: 404, Message: "Can't find thread author by nickname: " + thread.Author}
	}

	// todo mb limit 1?
	// todo insert consistency?
	errr := f.DB.QueryRow(`
	insert into threads 
    (title, author, message, forum, slug, created) 
	values ($1,$2,$3,(select slug from forums where slug = $4),$5,$6) 
	returning id, forum`, // todo right return?
		thread.Title, thread.Author, thread.Message, slug, threadSlug, thread.Created).Scan(&thread.ID, &thread.Forum)

	dbErr, _ := errr.(pgx.PgError)
	if dbErr.Code == pgerrcode.ForeignKeyViolation || dbErr.Code == pgerrcode.NotNullViolation {
		fmt.Println("404")
		fmt.Println(errr)
		return models.Thread{}, models.Error{Code: 404, Message: "Can't find thread forum by slug: " + thread.Slug}

	}
	if errr != nil { // если такой форум уже еть
		fmt.Println("409")
		return models.Thread{}, models.Error{Code: 409}
	}

	fmt.Println(thread)
	return thread, models.Error{Code: 201}
}

func (f *forumRepo) GetUsers(slug string, params models.ParseParams) ([]models.User, models.Error) {
	var queryParametres []interface{}
	query := `
		select userNickname, fullname, about, email from forum_users 
        where forumSlug = $1 
	`
	queryParametres = append(queryParametres, slug)

	if params.Since != "" {
		fmt.Println("with since", params.Since)
		if params.Desc == "true" {
			query += " and userNickname < $2 "
		} else {
			query += " and userNickname > $2 "
		}
		queryParametres = append(queryParametres, params.Since)
	}

	if params.Desc == "false" {
		query += " order by userNickname "
	} else {
		query += " order by userNickname DESC "
	}

	if params.Limit != 0 {
		if params.Since == "" {
			query += "  LIMIT $2"
		} else {
			query += "  LIMIT $3"
		}
		queryParametres = append(queryParametres, params.Limit)
	}

	forumUsers, err := f.DB.Query(query, queryParametres...)
	fmt.Println("            get forum users error:", err)
	if err != nil {
		return []models.User{}, models.Error{Code: 404}
	}

	users := make([]models.User, 0)
	for forumUsers.Next() {
		user := &models.User{}
		err = forumUsers.Scan(
			&user.Nickname,
			&user.FullName,
			&user.About,
			&user.Email,
		)
		if err != nil {
			fmt.Println("     scan forum user error:", err)
			return []models.User{}, models.Error{Code: 500}
		}

		users = append(users, *user)
	}

	return users, models.Error{}
}

func (f *forumRepo) GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error) {
	var queryParams []interface{}
	query := `
		select id, slug, forum, author, title, message, votes, created from threads 
		where forum = $1 `
	queryParams = append(queryParams, slug)

	if params.Since != "" {
		fmt.Println("with since")
		if params.Desc == "true" {
			query += "and created <= $2 "
		} else {
			query += "and created >= $2 "
		}

		queryParams = append(queryParams, params.Since)
	}

	if params.Desc == "false" {
		query += "order by created "
	} else {
		query += "order by created desc "
	}

	if params.Limit != 0 {
		if params.Since == "" {
			query += "LIMIT $2"
		} else {
			query += "LIMIT $3"
		}
		queryParams = append(queryParams, params.Limit)
	}

	forumThreads, err := f.DB.Query(query, queryParams...)
	fmt.Println(err)
	if err != nil {
		return []models.Thread{}, models.Error{Code: 404}
	}

	var threadSlug *string
	threads := make([]models.Thread, 0)
	for forumThreads.Next() {
		thread := new(models.Thread)
		err = forumThreads.Scan(
			&thread.ID,
			&threadSlug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
		)
		if threadSlug != nil {
			thread.Slug = *threadSlug
		}

		if err != nil {
			return []models.Thread{}, models.Error{Code: 404}
		}

		threads = append(threads, *thread)
	}

	if len(threads) == 0 {
		check := f.CheckForumExists(slug) // todo getBySlug???
		if !check {
			return threads, models.Error{Code: 404}
		}
	}

	return threads, models.Error{Code: 200}
}

func (f *forumRepo) GetThreadBySlug(slug string) (models.Thread, models.Error) {
	thread := models.Thread{}
	_ = f.DB.QueryRow("select slug, forum, author, title, message, votes, created  from threads where slug = $1 limit 1", slug).
		Scan(
			&thread.Slug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
			)
	return thread, models.Error{Code: 409}
}


func (f *forumRepo) CheckForumExists(slug string) (check bool) {
	_ = f.DB.QueryRow(`select exists(select 1 from forums where slug=$1)`, slug).Scan(&check)
	return check
}