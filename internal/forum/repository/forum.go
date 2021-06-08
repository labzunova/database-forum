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
	//
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
	fmt.Println("       get forum", slug)
	forumm := new(models.Forum)
	err := f.DB.QueryRow(`select title, "user", slug, posts_count, threads_count from forums where slug = $1`,
		slug).Scan(&forumm.Title, &forumm.User, &forumm.Slug, &forumm.Posts, &forumm.Threads)
	fmt.Println("forum", forumm)
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

	var id int
	err := f.DB.QueryRow("select id from users where nickname=$1", thread.Author).Scan(&id)
	if err == pgx.ErrNoRows {
		fmt.Println("404")
		return models.Thread{}, models.Error{Code: 404, Message: "Can't find thread author by nickname: " + thread.Author}
	}

	errr := f.DB.QueryRow(`
	insert into threads 
    (title, author, message, forum, slug, created) 
	values ($1,$2,$3,(select slug from forums where slug = $4),$5,$6) 
	returning id, forum`,
		thread.Title, thread.Author, thread.Message, slug, threadSlug, thread.Created).Scan(&thread.ID, &thread.Forum)

	dbErr, _ := errr.(pgx.PgError)
	if dbErr.Code == pgerrcode.ForeignKeyViolation || dbErr.Code == pgerrcode.NotNullViolation {
		fmt.Println("404")
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
		select userNickname, fullname, email, about from forum_users 
        where forumSlug = $1 
	`
	queryParametres = append(queryParametres, slug)

	if params.Since != "" {
		fmt.Println("with since", params.Since)
		if params.Desc {
			query += " and userNickname < $2 "
		} else {
			query += " and userNickname > $2 "
		}
		queryParametres = append(queryParametres, params.Since)
	}

	if !params.Desc {
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
		fmt.Println("start scan user")
		user := new(models.User)
		err = forumUsers.Scan(
			&user.Nickname,
			&user.FullName,
			&user.Email,
			&user.About,
		)
		fmt.Println("                uABOUT", user.About)
		fmt.Println(user)
		fmt.Println("     scan forum user error:", err)
		if err != nil {
			return []models.User{}, models.Error{Code: 500}
		}

		users = append(users, *user)
	}

	return users, models.Error{}
}

func (f *forumRepo) GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error) {
	fmt.Println("       get threads: ", slug, params)
fmt.Println("            ", slug)
	var queryParams []interface{}
	query := `
		select id, title, author, message, votes, slug, created from threads 
		where forum = $1 `
	queryParams = append(queryParams, slug)
	fmt.Println("            ", slug)

	if params.Since != "" {
		fmt.Println("with since")
		if params.Desc {
			query += "and created <= $2 "
		} else {
			query += "and created >= $2 "
		}

		//layout := "2006-01-02T15:04:05.000Z"
		//str := params.Since
		//t, _ := time.Parse(layout, str)

		queryParams = append(queryParams, params.Since)
	}
	fmt.Println("            ", slug)

	if !params.Desc {
		query += "order by created "
	} else {
		fmt.Println("with desc")
		query += "order by created desc "
	}
	fmt.Println("            ", slug)

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
	fmt.Println("            ", slug)

	forumThreads, err := f.DB.Query(query, queryParams...)
	fmt.Println("            ", slug)

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
			&thread.Title,
			&thread.Author,
			&thread.Message,
			&thread.Votes,
			&threadSlug,
			&thread.Created,
		)
		fmt.Println("            ", slug)
		fmt.Println(thread)
		thread.Forum = slug
		if threadSlug != nil {
			thread.Slug = *threadSlug
		}

		if err != nil {
			return []models.Thread{}, models.Error{Code: 404}
		}

		threads = append(threads, *thread)
	}

	fmt.Println(threads)
	if len(threads) == 0 {
		return threads, models.Error{Code: 404}
	}
	fmt.Println("            ", slug)
	fmt.Println("return")
	return threads, models.Error{Code: 200}
}

func (f *forumRepo) GetThreadBySlug(slug string) (models.Thread, models.Error) {
	thread := models.Thread{}
	_ = f.DB.QueryRow("select id, title, author, message, votes, slug, created, forum from threads where slug = $1", slug).
		Scan(&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
			&thread.Forum)
	return thread, models.Error{Code: 409}
}
