package repository

import (
	"DBproject/internal/forum"
	"DBproject/models"
	"database/sql"
	"fmt"
	"time"
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
	fmt.Println("create forum", forum)

	err := f.DB.QueryRow("select nickname from users where nickname=$1", forum.User).Scan(&forum.User)
	if err != nil {
		return models.Forum{}, models.Error{Code: 404}
	}

	// checkUserExists
	_, err = f.DB.Exec(`insert into forums 
    	(title, "user", slug) 
    	values ($1, $2, $3) returning "user"`,
		forum.Title, forum.User, forum.Slug)
	fmt.Println("create err ", err)
	//dbErr, ok := err.(pgx.PgError)
	//fmt.Println("create err ", dbErr, ok, dbErr.Code, dbErr.Message)
	//
	//if ok {
	//		if dbErr.Code == pgerrcode.NotNullViolation || dbErr.Code == pgerrcode.ForeignKeyViolation { // если владелец не найден
	//		return models.Forum{}, models.Error{Code: 404}
	//	}
	//	if dbErr.Code == pgerrcode.UniqueViolation { // если такой форум уже еть
	//		return models.Forum{}, models.Error{Code: 409}
	//	}
	//}

	if err != nil && err != sql.ErrNoRows { // если такой форум уже еть
		fmt.Println("409")
		return models.Forum{}, models.Error{Code: 409}
	}
	if err == sql.ErrNoRows { // если владелец не найден
		fmt.Println("404")
		return models.Forum{}, models.Error{Code: 404}
	}

	return forum, models.Error{Code: 200}
}

func (f *forumRepo) GetForum(slug string) (models.Forum, models.Error) {
	fmt.Println("get forum", slug)

	forum := new(models.Forum)
	err := f.DB.QueryRow(`select title, "user", slug, posts, threads from forums where slug = $1`,
		slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, models.Error{Code: 404}
	}

	return *forum, models.Error{Code: 200}
}

func (f *forumRepo) CreateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	fmt.Println("create thread: ", slug, thread)

	err := f.DB.QueryRow(
		`
	insert into threads 
    (title, author, message, forum, slug, created) 
	values ($1,$2,$3,$4,$5,$6) 
	returning id`,
		thread.Title, thread.Author, thread.Message, slug, thread.Slug, thread.Created).Scan(&thread.ID)

	//DBerror, _ := err.(pgx.PgError) // TODO error handling
	//switch DBerror.Code {
	//case pgerrcode.UniqueViolation: // если такой тред уже еть
	//	return thread, models.Error{Code: 409}
	//case pgerrcode.NotNullViolation: // если владелец не найден ???
	//	return models.Thread{}, models.Error{Code: 404}
	//}
fmt.Println(err)
	if err != nil && err != sql.ErrNoRows { // если такой форум уже еть
		fmt.Println("409")
		return models.Thread{}, models.Error{Code: 409}
	}
	if err == sql.ErrNoRows { // если владелец не найден
		fmt.Println("404")
		return models.Thread{}, models.Error{Code: 404}
	}
	    fmt.Println(thread)
	return thread, models.Error{Code: 201}
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
			&user.Nickname,
			&user.FullName,
			&user.Email,
			&user.About,
		)
		if err != nil {
			return []models.User{}, models.Error{Code: 500}
		}

		users = append(users, *user)
	}

	return users, models.Error{}
}

func (f *forumRepo) GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error) {
	fmt.Println("get threads: ", slug, params)

	var queryParams []interface{}
	query := `
		select id, title, author, message, votes, slug, created, forum from threads 
		where forum = $1 `
	queryParams = append(queryParams, slug)

	if params.Since != "" {
		fmt.Println("with since")
		if params.Desc {
			query += "and created <= $2 "
		} else {
			query += "and created >= $2 "
		}

		layout := "2006-01-02T15:04:05.000Z"
		str := params.Since
		t, _ := time.Parse(layout, str)
		t = t.Add(time.Hour * 3) // TODO ВРЕМЕННО ДЛЯ КОМПА

		queryParams = append(queryParams,t)
	}

	if !params.Desc {
		query += "order by created "
	} else {
		fmt.Println("with desc")
		query += "order by created desc "
	}

	if params.Limit != 0 {
		if params.Since == "" {
			fmt.Println("with limit1")
			query += "LIMIT $2"
		} else {
			fmt.Println("with limit2")
			query += "LIMIT $3"
		}
		queryParams = append(queryParams, params.Limit)
	}

	forumThreads, err := f.DB.Query(query, queryParams...)
	if err != nil {
		return []models.Thread{}, models.Error{Code: 404}
	}

	threads := make([]models.Thread, 0)
	for forumThreads.Next() {

		thread := new(models.Thread)
		err = forumThreads.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
			&thread.Forum,
		)

		thread.Created = thread.Created.Add(-time.Hour * 3) // TODO ВРЕМЕННО ДЛЯ КОМПА

		if err != nil {
			return  []models.Thread{}, models.Error{Code: 404}
		}

		threads = append(threads, *thread)
	}

	if len(threads) == 0 {
		return threads, models.Error{Code: 404}
	}
	return threads, models.Error{Code: 200}
}
