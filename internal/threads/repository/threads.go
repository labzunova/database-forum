package repository

import (
	"DBproject/internal/threads"
	"DBproject/models"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type threadsRepo struct {
	DB *pgx.ConnPool
}

func NewThreadsRepo(db *pgx.ConnPool) threads.ThreadsRepo {
	return &threadsRepo{
		DB: db,
	}
}

const treeDescSince = `
	select id, parent, author, message, isEdited, forum, thread, created
	from posts
	where thread = $1 and path < (select path from posts where id = $2)
	order by path desc
	limit $3`
const treeSince = `
	select id, parent, author, message, isEdited, forum, thread, created
	from posts
	where thread = $1 and path > (select path from posts where id = $2)
	order by path asc
	limit $3`
const treeDesc = `
	select id, parent, author, message, isEdited, forum, thread, created
	from posts
	where thread = $1 
	order by path desc
	limit $2`
const tree = `
	select id, parent, author, message, isEdited, forum, thread, created
	from posts
	where thread = $1 
	order by path asc
	limit $2`

const parentTreeDescSince = `
	select id, parent, author, message, isEdited, forum, thread, created
	from posts
	where thread = $1 and path[1] in(
		select distinct path[1] from posts
		where parent is null and thread = $1 and
		path[1]<(select path[1] from posts where id=$2)
		order by path[1] desc 
		limit $3
	)
	order by path[1] desc, path[2:]`
const parentTreeSince = `select id, parent, author, message, isEdited, forum, thread, created
	from posts
	where thread = $1 and path[1] in(
		select distinct path[1] from posts
		where parent is null and thread = $1 and
		path[1]>(select path[1] from posts where id=$2)
		order by path[1]
		limit $3
	)
	order by path`
const parentTreeDesc = `	
	select id, parent, author, message, isEdited, forum, thread, created
	from posts
	where thread = $1 and path[1] in(
		select distinct path[1] from posts
		where thread = $1
		order by path[1] desc 
		limit $2
	)
	order by path[1] desc, path[2:]`
const parentTree = `
	select id, parent, author, message, isEdited, forum, thread, created
	from posts
	where thread = $1 and path[1] in(
		select distinct path[1] from posts
		where thread = $1
		order by path[1]
		limit $2
	)
	order by path`


func (db *threadsRepo) GetThread(slug string, id int) (models.Thread, models.Error) {
	thread := models.Thread{
		Slug: slug,
		ID: id,
	}
	fmt.Println(slug)
	if id != 0 {
		fmt.Println("get thread by id")
		err := db.DB.QueryRow("select title, author, forum, message, votes, slug, created from threads where id = $1", id).
			Scan(&thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return thread, models.Error{Code: 404}
		}
	} else {
		fmt.Println("get thread by slug")
		err := db.DB.QueryRow("select id, slug, title, author, forum, message, votes, created from threads where slug = $1", slug).
			Scan(&thread.ID, &thread.Slug, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created)
		fmt.Println(thread, err)
		if err != nil {
			return thread, models.Error{Code: 404, Message: err.Error()}
		}
	}

	return thread, models.Error{Code: 200}
}

func (db *threadsRepo) UpdateThreadBySlug(slug string, thread models.Thread) (models.Thread, models.Error) {
	err := db.DB.QueryRow(`
		update threads set message=coalesce(nullif($1,''), message), title=coalesce(nullif($2,''), title) 
		where slug = $3 
		returning id, slug, title, author, forum, message, votes, created`,
		thread.Message, thread.Title, slug).
		Scan(&thread.ID, &thread.Slug, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return thread, models.Error{Code: 404, Message: err.Error()}
	}

	return thread, models.Error{Code: 200}
}

func (db *threadsRepo) UpdateThreadById(id int, thread models.Thread) (models.Thread, models.Error) {
	err := db.DB.QueryRow(`
		update threads set message=coalesce(nullif($1,''), message), title=coalesce(nullif($2,''), title)
		where id = $3 
		returning id, slug, title, author, forum, message, votes, created`,
		thread.Message, thread.Title, id).
		Scan(&thread.ID, &thread.Slug, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return thread, models.Error{Code: 404}
	}

	return thread, models.Error{Code: 200}
}

func (db *threadsRepo) GetThreadPostsById(id int, params models.ParseParamsThread) ([]models.Post, models.Error) {
	var query string
	var queryParameters []interface{}
	queryParameters = append(queryParameters, id)

	fmt.Println("since:", params.Since, "desc:", params.Desc, "limit:", params.Limit, "sort:", params.Sort)

	switch params.Sort {
	case "tree":
		switch {
		case params.Desc && params.Since != 0:
			query = treeDescSince
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case !params.Desc && params.Since != 0:
			query = treeSince
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case params.Desc && params.Since == 0:
			query = treeDesc
			queryParameters = append(queryParameters, params.Limit)
		case !params.Desc && params.Since == 0:
			query = tree
			queryParameters = append(queryParameters, params.Limit)
		default:
			return nil, models.Error{Code: 400}
		}

	case "parent_tree":
		switch {
		case params.Desc && params.Since != 0:
			query = parentTreeDescSince
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case !params.Desc && params.Since != 0:
			query = parentTreeSince
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case params.Desc && params.Since == 0:
			query = parentTreeDesc
			queryParameters = append(queryParameters, params.Limit)
		case !params.Desc && params.Since == 0:
			query = parentTree
			queryParameters = append(queryParameters, params.Limit)
		default:
			return nil, models.Error{Code: 400}
		}

	default: // flat
		switch {
		case params.Desc && params.Since != 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1 and id < $2
				ORDER BY id DESC
				LIMIT $3
			`
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case !params.Desc && params.Since != 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1 and id > $2
				ORDER BY id
				LIMIT $3
			`
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case params.Desc && params.Since == 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1
				ORDER BY id DESC
				LIMIT $2
			`
			queryParameters = append(queryParameters, params.Limit)
		case !params.Desc && params.Since == 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1
				ORDER BY id
				LIMIT $2
			`
			queryParameters = append(queryParameters, params.Limit)
		default:
			return nil, models.Error{Code: 400}
		}
	}

	rows, err := db.DB.Query(query, queryParameters...)
	fmt.Println("GET POSTS WITH SORT ERR ", err)
	if err == pgx.ErrNoRows {
		return nil, models.Error{Code: 404}
	}
	if err != nil {
		return nil, models.Error{Code: 500}
	}

	posts := make([]models.Post, 0)
	for rows.Next() {
		post := models.Post{}
		var parent *int
		err = rows.Scan(
			&post.ID,
			&parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if parent != nil {
			post.Parent = *parent
		}
		fmt.Println("POST: ", post)
		if err != nil {
			fmt.Println("error while scan", err)
			return nil, models.Error{Code: 500}
		}

		posts = append(posts, post)
	}
fmt.Println("enddd")
	return posts, models.Error{Code: 200}
}

func (db *threadsRepo) GetThreadPostsBySlug(slug string, params models.ParseParamsThread) ([]models.Post, models.Error) {
	return []models.Post{}, models.Error{} // todo
}

func (db *threadsRepo) VoteThreadBySlug(slug string, vote models.Vote) models.Error {
	_, err := db.DB.Exec(`INSERT INTO votes("user", thread, vote) values ($1,$2,$3)`, vote.Nickname, slug, vote.Voice)

	fmt.Println(err)
	dbErr, ok := err.(pgx.PgError)
	if ok {
		switch dbErr.Code {
		case pgerrcode.NotNullViolation:
			return models.Error{Code: 404}
		case pgerrcode.UniqueViolation:
			updateErr := db.UpdateVoteThreadBySlug(slug, vote)
			if updateErr.Code != 200 {
				return models.Error{Code: 500}
			}

			return models.Error{Code: 200}
		}
	}

	fmt.Println("NEW VOICE ", vote.Voice)
	_, err = db.DB.Exec("update threads set votes=votes+$1 where slug=$2", vote.Voice, slug)
	if err != nil {
		return models.Error{Code: 500}
	}

	return models.Error{Code: 200}
}

func (db *threadsRepo) UpdateVoteThreadBySlug(slug string, vote models.Vote) models.Error {
	fmt.Println("updatiiing")
	var voice int32
	err := db.DB.QueryRow(`select vote from votes 
	where thread=$1 and "user"=$2`, slug, vote.Nickname).Scan(&voice)
	fmt.Println(err)
	if err != nil {
		return models.Error{Code: 500}
	}
	fmt.Println("old voice ",voice)
	fmt.Println("new voice ",vote.Voice)

	if vote.Voice == voice {
		fmt.Println("RETURN")
		return models.Error{Code: 200}
	}

	_, err = db.DB.Exec(`UPDATE votes SET vote=$1 where thread=$2 and "user"=$3`, vote.Voice, slug, vote.Nickname)
	if err != nil {
		return models.Error{Code: 500}
	}

	if vote.Voice == -1 && voice == 1{
		vote.Voice = -2
	}

	if vote.Voice == 1 && voice == -1 {
		vote.Voice = 2
	}
	fmt.Println("UPDATING")

	fmt.Println("NEW VOICE ", vote.Voice)
	_, err = db.DB.Exec("update threads set votes=votes+$1 where slug=$2", vote.Voice, slug)
	if err != nil {
		return models.Error{Code: 500}
	}


	return models.Error{Code: 200}
}

func (db *threadsRepo) VoteThreadById(id int, vote models.Vote) models.Error {
	_, err := db.DB.Exec(`INSERT INTO votes("user", thread, vote) values 
    	($1,(select slug from threads where id=$2),$3)`, vote.Nickname, id, vote.Voice)
	fmt.Println(err)
	dbErr, ok := err.(pgx.PgError)
	fmt.Println(dbErr.Code)
	if ok {
		switch dbErr.Code {
		case pgerrcode.NotNullViolation:
			return models.Error{Code: 404}
		case pgerrcode.UniqueViolation:
			updateErr := db.UpdateVoteThreadById(id, vote)
			if updateErr.Code != 200 {
				return models.Error{Code: 500}
			}
			return models.Error{Code: 200}
		}
	}
	fmt.Println("UPDATING")

	fmt.Println("NEW VOICE ", vote)
	_, err = db.DB.Exec("update threads set votes=votes+$1 where id=$2", vote.Voice, id)
	if err != nil {
		return models.Error{Code: 500}
	}

	return models.Error{Code: 200}
}

func (db *threadsRepo) UpdateVoteThreadById(id int, vote models.Vote) models.Error {
	//err := db.DB.QueryRow(`UPDATE votes SET vote=(coalesce(nullif($1,vote), vote))
	//	where thread=(select slug from threads where id=$2) and user=$3
	//	returning vote`, vote.Voice, id, vote.Nickname).Scan(&vote.Voice)
	var voice int32
	fmt.Println("id: ", id, "user: ", vote.Nickname)
	err := db.DB.QueryRow(`select vote from votes 
	where thread=(select slug from threads where id=$1) and "user"=$2`, id, vote.Nickname).Scan(&voice)
	fmt.Println("old voice ",voice)
	fmt.Println("new voice ",vote.Voice)
	if err != nil {
		return models.Error{Code: 500}
	}

	if vote.Voice == voice {
		return models.Error{Code: 200}
	}

	_, err = db.DB.Exec(`UPDATE votes SET vote=$1 
		where thread=(select slug from threads where id=$2) and "user"=$3`, vote.Voice, id, vote.Nickname)
	fmt.Println(3453534,err)
	if err != nil {
		return models.Error{Code: 500}
	}

	if vote.Voice == -1 && voice == 1 {
		vote.Voice = -2
	}

	if vote.Voice == 1 && voice == -1 {
		vote.Voice = 2
	}

	fmt.Println("UPDATING")
	_, err = db.DB.Exec("update threads set votes=votes+$1 where id=$2", vote.Voice, id)
	fmt.Println(3453534,err)
	if err != nil {
		return models.Error{Code: 500}
	}

	return models.Error{Code: 200}
}

func (db *threadsRepo) GetThreadIDBySlug(slug string) (int, models.Error) {
	var id int
	_ = db.DB.QueryRow("select id from threads where slug=$1", slug).Scan(&id)
	return id, models.Error{Code: 200}
}
