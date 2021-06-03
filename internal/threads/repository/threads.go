package repository

import (
	"DBproject/internal/threads"
	"DBproject/models"
	"database/sql"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type threadsRepo struct {
	DB *sql.DB
}

func NewThreadsRepo(db *sql.DB) threads.ThreadsRepo {
	return &threadsRepo{
		DB: db,
	}
}

func (db *threadsRepo) GetThread(slug string, id int) (models.Thread, models.Error) {
	thread := models.Thread{
		Slug: slug,
		ID: id,
	}
	if id == 0 {
		err := db.DB.QueryRow("select title, author, forum, message, votes, slug, created from threads where id = $1", id).
			Scan(&thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return thread, models.Error{Code: 404}
		}
	} else {
		err := db.DB.QueryRow("select id, title, author, forum, message, votes, created from threads where slug = $1", slug).
			Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created)
		if err != nil {
			return thread, models.Error{Code: 404}
		}
	}

	return thread, models.Error{}
}

func (db *threadsRepo) UpdateThreadBySlug(slug string, thread models.Thread) (models.Thread, models.Error) {
	err := db.DB.QueryRow(`
		update threads set message=coalesce(nullif($1,""), message), title=coalesce(nullif($2,""), title)
		where slug = $3 
		returning id, title, author, forum, message, votes, created`,
		thread.Title, thread.Message, slug).
		Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return thread, models.Error{Code: 404}
	}

	return thread, models.Error{Code: 200}
}

func (db *threadsRepo) UpdateThreadById(id int, thread models.Thread) (models.Thread, models.Error) {
	err := db.DB.QueryRow(`
		update threads set message=coalesce(nullif($1,""), message), title=coalesce(nullif($2,""), title)
		where id = $3 
		returning id, title, author, forum, message, votes, created`,
		thread.Title, thread.Message, id).
		Scan(&thread.ID, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Created)
	if err != nil {
		return thread, models.Error{Code: 404}
	}

	return thread, models.Error{Code: 200}
}

func (db *threadsRepo) GetThreadPostsById(id int, params models.ParseParamsThread) ([]models.Post, models.Error) {
	var query string
	var queryParameters []interface{}
	queryParameters = append(queryParameters, id)

	switch params.Sort {
	case "tree":
		switch {
		case params.Desc && params.Since != 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1 and id
`
		case !params.Desc && params.Since != 0:

		case params.Desc && params.Since == 0:

		case !params.Desc && params.Since == 0:

		default:
			return nil, models.Error{Code: 400}
		}

	case "parent_tree":
		switch {
		case params.Desc && params.Since != 0:

		case !params.Desc && params.Since != 0:

		case params.Desc && params.Since == 0:

		case !params.Desc && params.Since == 0:

		default:
			return nil, models.Error{Code: 400}
		}

	default: // flat
		switch {
		case params.Desc && params.Since != 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1 and id > $2
				ORDER BY created DESC
				LIMIT $3
			`
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case !params.Desc && params.Since != 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1 and id > $2
				ORDER BY created
				LIMIT $2
			`
			queryParameters = append(queryParameters, params.Since, params.Limit)
		case params.Desc && params.Since == 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1
				ORDER BY created DESC
				LIMIT $2
			`
			queryParameters = append(queryParameters, params.Limit)
		case !params.Desc && params.Since == 0:
			query = `
				SELECT id, parent, author, message, isEdited, forum, thread, created
				FROM posts
				WHERE thread = $1
				ORDER BY created
				LIMIT $2
			`
			queryParameters = append(queryParameters, params.Limit)
		default:
			return nil, models.Error{Code: 400}
		}
	}

	posts := make([]models.Post, 0)
	rows, err := db.DB.Query(query, queryParameters)
	if err == sql.ErrNoRows {
		return nil, models.Error{Code: 404}
	}
	if err != nil {
		return nil, models.Error{Code: 500}
	}

	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(
			&post.ID,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			return nil, models.Error{Code: 500}
		}
	}

	return posts, models.Error{Code: 200}
}

func (db *threadsRepo) GetThreadPostsBySlug(slug string, params models.ParseParamsThread) ([]models.Post, models.Error) {
	return []models.Post{}, models.Error{} // todo
}

func (db *threadsRepo) VoteThreadBySlug(slug string, vote models.Vote) models.Error {
	_, err := db.DB.Exec("INSERT INTO votes(user, thread, vote) values ($1,$2,$3)", vote.Nickname, slug, vote.Voice)
	dbErr := err.(pgx.PgError)
	switch dbErr.Code {
	case pgerrcode.NotNullViolation:
		return models.Error{Code: 404}
	case pgerrcode.UniqueViolation:
		updateErr := db.UpdateVoteThreadBySlug(slug, vote)
		if updateErr.Code != 0 {
			return models.Error{Code: 500}
		}
		return models.Error{Code: 200}
	}

	_, err = db.DB.Exec("update threads set votes=votes+1 where slug=$1", slug)
	if err != nil {
		return models.Error{Code: 500}
	}

	return models.Error{Code: 200}
}

func (db *threadsRepo) UpdateVoteThreadBySlug(slug string, vote models.Vote) models.Error {
	_, err := db.DB.Exec("UPDATE votes SET vote=vote+$1 where thread=$2", vote.Voice, slug)
	if err != nil {
		return models.Error{Code: 500}
	}

	return models.Error{Code: 200}
}

func (db *threadsRepo) VoteThreadById(id int, vote models.Vote) models.Error {
	_, err := db.DB.Exec(`INSERT INTO votes(user, thread, vote) values 
    	($1,$2, (select slug from threads where id=$3))`, vote.Nickname, id, vote.Voice)
	dbErr := err.(pgx.PgError)
	switch dbErr.Code {
	case pgerrcode.NotNullViolation:
		return models.Error{Code: 404}
	case pgerrcode.UniqueViolation:
		updateErr := db.UpdateVoteThreadById(id, vote)
		if updateErr.Code != 0 {
			return models.Error{Code: 500}
		}
		return models.Error{Code: 200}
	}

	_, err = db.DB.Exec("update threads set votes=votes+1 where id=$1", id)
	if err != nil {
		return models.Error{Code: 500}
	}

	return models.Error{Code: 200}
}

func (db *threadsRepo) UpdateVoteThreadById(id int, vote models.Vote) models.Error {
	_, err := db.DB.Exec("UPDATE votes SET vote=vote+$1 where thread=(select slug from threads where id=$2)", vote.Voice, id)
	if err != nil {
		return models.Error{Code: 500}
	}

	return models.Error{Code: 200}
}
