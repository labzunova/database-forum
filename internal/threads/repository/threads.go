package repository

import (
	"DBproject/internal/threads"
	"DBproject/models"
	"fmt"
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
	select thread, author, forum, isEdited, message, parent, created, id
	from posts
	where thread = $1 and path < (select path from posts where id = $2)
	order by path desc
	limit $3`
const treeSince = `
	select thread, author, forum, isEdited, message, parent, created, id
	from posts
	where thread = $1 and path > (select path from posts where id = $2)
	order by path asc
	limit $3`
const treeDesc = `
	select thread, author, forum, isEdited, message, parent, created, id
	from posts
	where thread = $1 
	order by path desc
	limit $2`
const tree = `
	select thread, author, forum, isEdited, message, parent, created, id
	from posts
	where thread = $1 
	order by path asc
	limit $2`

const parentTreeDescSince = `
	select thread, author, forum, isEdited, message, parent, created, id
	from posts
	where thread = $1 and path[1] in(
		select distinct path[1] from posts
		where parent is null and thread = $1 and
		path[1]<(select path[1] from posts where id=$2)
		order by path[1] desc 
		limit $3
	)
	order by path[1] desc, path[2:]`
const parentTreeSince = `
	select thread, author, forum, isEdited, message, parent, created, id
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
	select thread, author, forum, isEdited, message, parent, created, id
	from posts
	where thread = $1 and path[1] in(
		select distinct path[1] from posts
		where thread = $1
		order by path[1] desc 
		limit $2
	)
	order by path[1] desc, path[2:]`
const parentTree = `
	select thread, author, forum, isEdited, message, parent, created, id
	from posts
	where thread = $1 and path[1] in(
		select distinct path[1] from posts
		where thread = $1
		order by path[1]
		limit $2
	)
	order by path`

func (db *threadsRepo) GetThread(slug string, id int) (models.Thread, models.Error) {
	var threadSlug *string
	thread := models.Thread{}
	if id != 0 {
		thread.ID = id
	} else {
		thread.Slug = slug
	}

	if id != 0 {
		err := db.DB.QueryRow("select id, slug, forum, author, title, message, votes, created from threads where id = $1 limit 1", id).
			Scan(&thread.ID, &threadSlug, &thread.Forum, &thread.Author, &thread.Title, &thread.Message, &thread.Votes, &thread.Created)
		fmt.Println(err)
		if threadSlug != nil {
			thread.Slug = *threadSlug
		}
		if err != nil {
			return thread, models.Error{Code: 404}
		}
		return thread, models.Error{Code: 200}
	}

	// by slug
		err := db.DB.QueryRow("select id, slug, forum, author, title, message, votes, created from threads where slug = $1", slug).
			Scan(&thread.ID, &threadSlug, &thread.Forum, &thread.Author, &thread.Title, &thread.Message, &thread.Votes, &thread.Created)
		if threadSlug != nil {
			thread.Slug = *threadSlug
		}
		fmt.Println(err)
		if err != nil {
			fmt.Println("Ererer")
			return thread, models.Error{Code: 404, Message: "error"}
		}
		return thread, models.Error{Code: 200}

}

func (db *threadsRepo) UpdateThreadBySlug(slug string, thread models.Thread) (models.Thread, models.Error) {
	var threadSlug *string
	err := db.DB.QueryRow(`
		update threads set message=coalesce(nullif($1,''), message), title=coalesce(nullif($2,''), title) 
		where slug = $3 
		returning id, slug, forum, author, title, message, votes, created`,
		thread.Message, thread.Title, slug).
		Scan(&thread.ID, &threadSlug, &thread.Forum, &thread.Author, &thread.Title, &thread.Message, &thread.Votes, &thread.Created)
	if threadSlug != nil {
		thread.Slug = *threadSlug
	}
	if err != nil {
		return thread, models.Error{Code: 404, Message: err.Error()}
	}

	return thread, models.Error{Code: 200}
}

func (db *threadsRepo) UpdateThreadById(id int, thread models.Thread) (models.Thread, models.Error) {
	var threadSlug *string
	err := db.DB.QueryRow(`
		update threads set message=coalesce(nullif($1,''), message), title=coalesce(nullif($2,''), title)
		where id = $3 
		returning id, slug, forum, author, title, message, votes, created`,
		thread.Message, thread.Title, id).
		Scan(&thread.ID, &threadSlug, &thread.Forum, &thread.Author, &thread.Title, &thread.Message, &thread.Votes, &thread.Created)
	if threadSlug != nil {
		thread.Slug = *threadSlug
	}
	if err != nil {
		return thread, models.Error{Code: 404}
	}

	return thread, models.Error{Code: 200}
}

func (db *threadsRepo) GetThreadPostsById(id int, slugOrId string, params models.ParseParamsThread) ([]models.Post, models.Error) {
	var query string
	var queryParameters []interface{}
	queryParameters = append(queryParameters, id)

	if params.Limit == 0 {
		params.Limit = 1000000
	}

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
				select thread, author, forum, isEdited, message, parent, created, id
				FROM posts
				WHERE thread = $1 and id < $2
				ORDER BY id DESC
				LIMIT $3
			`
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case !params.Desc && params.Since != 0:
			query = `
				select thread, author, forum, isEdited, message, parent, created, id
				FROM posts
				WHERE thread = $1 and id > $2
				ORDER BY id
				LIMIT $3
			`
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case params.Desc && params.Since == 0:
			query = `
				select thread, author, forum, isEdited, message, parent, created, id
				FROM posts
				WHERE thread = $1
				ORDER BY id DESC
				LIMIT $2
			`
			queryParameters = append(queryParameters, params.Limit)
		case !params.Desc && params.Since == 0:
			query = `
				select thread, author, forum, isEdited, message, parent, created, id
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
	//if err == pgx.ErrNoRows {
	//	_, ok := strconv.Atoi(slugOrId)
	//	if ok == nil {
	//		return nil, models.Error{Code: 404, Message: fmt.Sprintf("Can't find forum by id: %d", id)}
	//	}
	//	return nil, models.Error{Code: 404, Message: "Can't find forum by slug: " + slugOrId}
	//}
	if err != nil {
		return nil, models.Error{Code: 404, Message: "Can't find forum "}
	}

	posts := make([]models.Post, 0)
	for rows.Next() {
		post := models.Post{}
		var parent *int
		err = rows.Scan(
			&post.Thread,
			&post.Author,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&parent,
			&post.Created,
			&post.ID,
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
	fmt.Println(posts)
	fmt.Println(id)
	fmt.Println("enddd")
	return posts, models.Error{Code: 200}
}

func (db *threadsRepo) GetThreadPostsBySlug(slug string, params models.ParseParamsThread) ([]models.Post, models.Error) {
	return []models.Post{}, models.Error{} // todo
}

func (db *threadsRepo) VoteThreadBySlug(slug string, vote models.Vote) models.Error {
	_, err := db.DB.Exec(`INSERT INTO votes("user", thread, vote) 
		values ($1,(select id from threads where slug=$2),$3)
		on conflict ("user",thread) do
		update set vote=$3`, vote.Nickname, slug, vote.Voice)

	fmt.Println(err)
	if err != nil {
		fmt.Println("ERROR")
		return models.Error{Code: 404}
	}

	return models.Error{Code: 200}
}

func (db *threadsRepo) VoteThreadById(id int, vote models.Vote) models.Error {
	_, err := db.DB.Exec(`INSERT INTO votes("user", thread, vote) values 
    	($1,$2,$3)
		on conflict ("user",thread) do
		update set vote=$3`, vote.Nickname, id, vote.Voice)
	if err != nil {
		return models.Error{Code: 404}
	}
	fmt.Println("UPDATING")

	return models.Error{Code: 200}
}

func (db *threadsRepo) GetThreadIDBySlug(slug string, id int) (int, models.Error) {
	err := db.DB.QueryRow("select id from threads where slug=$1 or id=$2", slug, id).Scan(&id)
	if err != nil || id == 0 {
		return id, models.Error{Code: 404}
	}
	return id, models.Error{Code: 200}
}
